/*
Copyright © 2020 Marvin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package csv

import (
	"fmt"
	"strings"
	"time"

	"github.com/wentaojin/transferdb/utils"

	"go.uber.org/zap"

	"github.com/wentaojin/transferdb/service"
	"github.com/xxjwxc/gowp/workpool"
)

func startOracleTableFullCSV(cfg *service.CfgFile, engine *service.Engine, waitSyncTableInfo, partSyncTableInfo []string, syncMode string) error {
	characterSet, err := engine.GetOracleDBCharacterSet()
	if err != nil {
		return err
	}
	isGBKCharacterSet := false
	if strings.Contains(strings.ToUpper(characterSet), ".ZHS16GBK") {
		isGBKCharacterSet = true
	}
	var OracleCharacterSet string
	if isGBKCharacterSet {
		OracleCharacterSet = utils.OracleUTF8CharacterSet
	} else {
		OracleCharacterSet = utils.OracleGBKCharacterSet
	}
	if len(partSyncTableInfo) > 0 {
		if err := startOracleTableConsumeByCheckpoint(cfg, engine, partSyncTableInfo, OracleCharacterSet, syncMode); err != nil {
			return err
		}
	}
	if len(waitSyncTableInfo) > 0 {
		if err := startOracleTableConsumeBySCN(cfg, engine, waitSyncTableInfo, OracleCharacterSet, syncMode); err != nil {
			return err
		}
	}
	return nil
}

func startOracleTableConsumeByCheckpoint(cfg *service.CfgFile, engine *service.Engine, partSyncTableInfo []string, sourceCharset, syncMode string) error {
	wp := workpool.New(cfg.CSVConfig.WorkerThreads)

	for _, tbl := range partSyncTableInfo {
		table := tbl
		wp.Do(func() error {
			if err := syncOracleRowsByRowID(cfg, engine, sourceCharset, table, syncMode); err != nil {
				return err
			}
			return nil
		})
	}
	if err := wp.Wait(); err != nil {
		return err
	}
	if !wp.IsDone() {
		return fmt.Errorf("sync oracle table rows by checkpoint failed, please rerunning")
	}
	return nil
}

func startOracleTableConsumeBySCN(cfg *service.CfgFile, engine *service.Engine, waitSyncTableInfo []string, sourceCharset, syncMode string) error {
	wp := workpool.New(cfg.CSVConfig.WorkerThreads)

	for idx, tbl := range waitSyncTableInfo {
		table := tbl
		seq := idx
		wp.Do(func() error {
			startTime := time.Now()
			service.Logger.Info("single full table init scn start",
				zap.String("schema", cfg.SourceConfig.SchemaName),
				zap.String("table", table))

			// 全量同步前，获取 SCN 以及初始化元数据表
			globalSCN, err := engine.GetOracleCurrentSnapshotSCN()
			if err != nil {
				return err
			}
			if err = engine.InitWaitAndFullSyncMetaRecord(cfg.SourceConfig.SchemaName,
				table, seq, globalSCN, cfg.CSVConfig.ChunkSize, cfg.AppConfig.InsertBatchSize, syncMode); err != nil {
				return err
			}

			endTime := time.Now()
			service.Logger.Info("single full table init scn finished",
				zap.String("schema", cfg.SourceConfig.SchemaName),
				zap.String("table", table),
				zap.String("cost", endTime.Sub(startTime).String()))

			if err = syncOracleRowsByRowID(cfg, engine, sourceCharset, table, syncMode); err != nil {
				return err
			}
			return nil
		})
	}
	if err := wp.Wait(); err != nil {
		return err
	}
	if !wp.IsDone() {
		return fmt.Errorf("sync oracle table rows by scn failed, please rerunning")
	}
	return nil
}

func syncOracleRowsByRowID(cfg *service.CfgFile, engine *service.Engine, sourceCharset, sourceTableName, syncMode string) error {
	startTime := time.Now()
	service.Logger.Info("single full table data sync start",
		zap.String("schema", cfg.SourceConfig.SchemaName),
		zap.String("charset", sourceCharset),
		zap.String("table", sourceTableName))

	oraRowIDSQL, err := engine.GetFullSyncMetaRowIDRecord(cfg.SourceConfig.SchemaName, sourceTableName)
	if err != nil {
		return err
	}
	wp := workpool.New(cfg.CSVConfig.TableThreads)
	for idx, rowidSQL := range oraRowIDSQL {
		sql := rowidSQL
		dirIndex := idx
		wp.DoWait(func() error {
			// 抽取 Oracle 数据
			var (
				columnFields []string
				rowsResult   [][]string
			)
			columnFields, rowsResult, err = extractorTableFullRecord(engine, cfg.CSVConfig, cfg.SourceConfig.SchemaName, sourceTableName, sql)
			if err != nil {
				return err
			}

			if len(rowsResult) == 0 {
				service.Logger.Warn("oracle schema table rowid data return null rows, skip",
					zap.String("schema", cfg.SourceConfig.SchemaName),
					zap.String("table", sourceTableName),
					zap.String("sql", sql))
				// 清理记录以及更新记录
				if err = engine.ModifyWaitAndFullSyncTableMetaRecord(
					cfg.TargetConfig.MetaSchema,
					cfg.SourceConfig.SchemaName, sourceTableName, sql, syncMode); err != nil {
					return err
				}
				return nil
			}

			// 转换/应用 Oracle 数据 -> MySQL
			if err = applierTableFullRecord(cfg.TargetConfig.SchemaName,
				sourceTableName, sql, cfg.CSVConfig.ApplyThreads,
				translatorTableFullRecord(cfg.TargetConfig.SchemaName, sourceTableName,
					sql, sourceCharset, dirIndex, columnFields, rowsResult, cfg.CSVConfig)); err != nil {
				return err
			}

			// 清理记录以及更新记录
			if err = engine.ModifyWaitAndFullSyncTableMetaRecord(
				cfg.TargetConfig.MetaSchema,
				cfg.SourceConfig.SchemaName, sourceTableName, sql, syncMode); err != nil {
				return err
			}
			return nil
		})
	}
	if err = wp.Wait(); err != nil {
		return err
	}

	endTime := time.Now()
	if !wp.IsDone() {
		service.Logger.Fatal("single full table data loader failed",
			zap.String("schema", cfg.SourceConfig.SchemaName),
			zap.String("table", sourceTableName),
			zap.String("cost", endTime.Sub(startTime).String()))
		return fmt.Errorf("oracle schema [%s] single full table [%v] data loader failed",
			cfg.SourceConfig.SchemaName, sourceTableName)
	}
	service.Logger.Info("single full table data loader finished",
		zap.String("schema", cfg.SourceConfig.SchemaName),
		zap.String("table", sourceTableName),
		zap.String("cost", endTime.Sub(startTime).String()))

	return nil
}
