package reverser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wentaojin/transferdb/utils"

	"github.com/wentaojin/transferdb/service"

	"go.uber.org/zap"
)

func ReverseOracleToMySQLTableV2(engine *service.Engine, cfg *service.CfgFile) ([]string, error) {
	files := make([]string, 0)
	startTime := time.Now()
	service.Logger.Info("reverse table oracle to mysql start",
		zap.String("schema", cfg.SourceConfig.SchemaName))

	// 用于检查下游 MySQL/TiDB 环境检查
	// 只提供表结构转换文本输出，不提供直写下游，故注释下游检查项
	//if err := reverseOracleToMySQLTableInspect(engine, cfg); err != nil {
	//	return err
	//}

	// 获取待转换表
	exporterTableSlice, err := cfg.GenerateTables(engine)
	if err != nil {
		return nil, err
	}

	if len(exporterTableSlice) == 0 {
		service.Logger.Warn("there are no table objects in the oracle schema",
			zap.String("schema", cfg.SourceConfig.SchemaName))
		return nil, fmt.Errorf("there are no table objects in the oracle schema")
	}

	// 判断 table_error_detail 是否存在错误记录，是否可进行 reverse
	errorTotals, err := engine.GetTableErrorDetailCountByMode(cfg.SourceConfig.SchemaName, utils.ReverseMode)
	if err != nil {
		return nil, fmt.Errorf("func [GetTableErrorDetailCountByMode] reverse schema [%s] table mode [%s] task failed, error: %v", strings.ToUpper(cfg.SourceConfig.SchemaName), utils.ReverseMode, err)
	}
	if errorTotals > 0 {
		return nil, fmt.Errorf("func [GetTableErrorDetailCountByMode] reverse schema [%s] table mode [%s] task failed, table [table_error_detail] exist failed error, please clear and rerunning", strings.ToUpper(cfg.SourceConfig.SchemaName), utils.ReverseMode)
	}

	// oracle db collation
	nlsSort, err := engine.GetOracleDBCharacterNLSSortCollation()
	if err != nil {
		return nil, err
	}
	nlsComp, err := engine.GetOracleDBCharacterNLSCompCollation()
	if err != nil {
		return nil, err
	}
	if _, ok := utils.OracleCollationMap[strings.ToUpper(nlsSort)]; !ok {
		return nil, fmt.Errorf("oracle db nls sort [%s] isn't support", nlsSort)
	}
	if _, ok := utils.OracleCollationMap[strings.ToUpper(nlsComp)]; !ok {
		return nil, fmt.Errorf("oracle db nls comp [%s] isn't support", nlsComp)
	}
	if strings.ToUpper(nlsSort) != strings.ToUpper(nlsComp) {
		return nil, fmt.Errorf("oracle db nls_sort [%s] and nls_comp [%s] isn't different, need be equal; because mysql db isn't support", nlsSort, nlsComp)
	}

	// 表列表
	tables, partitionTableList, temporaryTableList, clusteredTableList, err := LoadOracleToMySQLTableList(engine, exporterTableSlice, cfg.SourceConfig.SchemaName, cfg.TargetConfig.SchemaName, nlsSort, nlsComp, cfg.TargetConfig.Overwrite, cfg.AppConfig.Threads)
	if err != nil {
		return nil, err
	}

	var (
		pwdDir                         string
		fileReverse, fileCompatibility *os.File
	)
	pwdDir, err = os.Getwd() //这里是获取本地路径做为
	if err != nil {
		return nil, err
	}

	fileReverse, err = os.OpenFile(filepath.Join(pwdDir, fmt.Sprintf("reverse_%s.sql", cfg.SourceConfig.SchemaName)), os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer fileReverse.Close()

	fileCompatibility, err = os.OpenFile(filepath.Join(pwdDir, fmt.Sprintf("compatibility_%s.sql", cfg.SourceConfig.SchemaName)), os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer fileCompatibility.Close()

	// 需要返回的数据
	files = append(files, fileReverse.Name(), fileCompatibility.Name())

	wrReverse := &FileMW{sync.Mutex{}, fileReverse}
	wrComp := &FileMW{sync.Mutex{}, fileCompatibility}

	// 创建 Schema
	if err := GenCreateSchema(wrReverse, engine, strings.ToUpper(cfg.SourceConfig.SchemaName), strings.ToUpper(cfg.TargetConfig.SchemaName), nlsComp); err != nil {
		return nil, err
	}

	// 不兼容项 - 表提示
	if err = CompatibilityDBTips(wrComp, strings.ToUpper(cfg.SourceConfig.SchemaName), partitionTableList, temporaryTableList, clusteredTableList); err != nil {
		return nil, err
	}

	// 设置工作池
	// 设置 goroutine 数
	wg := sync.WaitGroup{}
	ch := make(chan Table, utils.BufferSize)

	go func() {
		for _, t := range tables {
			ch <- t
		}
		close(ch)
	}()

	for c := 0; c < cfg.AppConfig.Threads; c++ {
		wg.Add(1)
		go func(revFileMW, compFileMW *FileMW) {
			defer wg.Done()
			for t := range ch {
				writer, err := NewReverseWriter(t, revFileMW, compFileMW)
				if err != nil {
					if err = t.Engine.GormDB.Create(&service.TableErrorDetail{
						SourceSchemaName: t.SourceSchemaName,
						SourceTableName:  t.SourceTableName,
						RunMode:          utils.ReverseMode,
						InfoSources:      utils.ReverseMode,
						RunStatus:        "Failed",
						Detail:           t.String(),
						Error:            err.Error(),
					}).Error; err != nil {
						service.Logger.Error("reverse table oracle to mysql failed",
							zap.String("schema", t.SourceSchemaName),
							zap.String("table", t.SourceTableName),
							zap.Error(
								fmt.Errorf("func [NewReverseWriter] reverse table task failed, detail see [table_error_detail], please rerunning")))
						panic(
							fmt.Errorf("func [NewReverseWriter] reverse table task failed, detail see [table_error_detail], please rerunning, error: %v", err))
					}
					continue
				}
				if err = writer.Reverse(); err != nil {
					if err = t.Engine.GormDB.Create(&service.TableErrorDetail{
						SourceSchemaName: t.SourceSchemaName,
						SourceTableName:  t.SourceTableName,
						RunMode:          utils.ReverseMode,
						InfoSources:      utils.ReverseMode,
						RunStatus:        "Failed",
						Detail:           t.String(),
						Error:            err.Error(),
					}).Error; err != nil {
						service.Logger.Error("reverse table oracle to mysql failed",
							zap.String("scheme", t.SourceSchemaName),
							zap.String("table", t.SourceTableName),
							zap.Error(
								fmt.Errorf("func [Reverse] reverse table task failed, detail see [table_error_detail], please rerunning")))
						panic(
							fmt.Errorf("func [Reverse] reverse table task failed, detail see [table_error_detail], please rerunning, error: %v", err))
					}
					continue
				}
			}
		}(wrReverse, wrComp)
	}

	wg.Wait()

	errorTotals, err = engine.GetTableErrorDetailCountByMode(cfg.SourceConfig.SchemaName, utils.ReverseMode)
	if err != nil {
		return nil, fmt.Errorf("func [GetTableErrorDetailCountByMode] reverse schema [%s] mode [%s] table task failed, error: %v", strings.ToUpper(cfg.SourceConfig.SchemaName), utils.ReverseMode, err)
	}

	endTime := time.Now()
	service.Logger.Info("reverse", zap.String("create table and index output", filepath.Join(pwdDir,
		fmt.Sprintf("reverse_%s.sql", cfg.SourceConfig.SchemaName))))
	service.Logger.Info("compatibility", zap.String("maybe exist compatibility output", filepath.Join(pwdDir,
		fmt.Sprintf("compatibility_%s.sql", cfg.SourceConfig.SchemaName))))
	if errorTotals == 0 {
		service.Logger.Info("reverse table oracle to mysql finished",
			zap.Int("table totals", len(tables)),
			zap.Int("table success", len(tables)),
			zap.Int("table failed", int(errorTotals)),
			zap.String("cost", endTime.Sub(startTime).String()))
	} else {
		service.Logger.Warn("reverse table oracle to mysql finished",
			zap.Int("table totals", len(tables)),
			zap.Int("table success", len(tables)-int(errorTotals)),
			zap.Int("table failed", int(errorTotals)),
			zap.String("failed tips", "failed detail, please see table [table_error_detail]"),
			zap.String("cost", endTime.Sub(startTime).String()))
	}
	return files, nil
}
