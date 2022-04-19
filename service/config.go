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
package service

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/BurntSushi/toml"
)

// 程序配置文件
//todo 需增加一个sync包， 保证并发， 后面再做
type CfgFile struct {
	AppConfig    AppConfig    `toml:"app" json:"app"`
	FullConfig   FullConfig   `toml:"full" json:"full"`
	CSVConfig    CSVConfig    `toml:"csv" json:"csv"`
	AllConfig    AllConfig    `toml:"all" json:"all"`
	SourceConfig SourceConfig `toml:"source" json:"source"`
	TargetConfig TargetConfig `toml:"target" json:"target"`
	LogConfig    LogConfig    `toml:"log" json:"log"`
	DiffConfig   DiffConfig   `toml:"diff" json:"diff"`
}

type AppConfig struct {
	InsertBatchSize  int    `toml:"insert-batch-size" json:"insert-batch-size"`
	SlowlogThreshold int    `toml:"slowlog-threshold" json:"slowlog-threshold"`
	Threads          int    `toml:"threads" json:"threads"`
	PprofPort        string `toml:"pprof-port" json:"pprof-port"`
	HTTPPort         string `toml:"http-port" json:"http-port"`
	JwtSecret        string `toml:"jwt-secret" json:"jwt-secret"`
}

type DiffConfig struct {
	ChunkSize         int           `toml:"chunk-size" json:"chunk-size"`
	DiffThreads       int           `toml:"diff-threads" json:"diff-threads"`
	OnlyCheckRows     bool          `toml:"only-check-rows" json:"only-check-rows"`
	EnableCheckpoint  bool          `toml:"enable-checkpoint" json:"enable-checkpoint"`
	IgnoreStructCheck bool          `toml:"ignore-struct-check" json:"ignore-struct-check"`
	FixSqlFile        string        `toml:"fix-sql-file" json:"fix-sql-file"`
	TableConfig       []TableConfig `toml:"table-config" json:"table-config"`
}

type TableConfig struct {
	SourceTable string `toml:"source-table" json:"source-table"`
	IndexFields string `toml:"index-fields" json:"index-fields"`
	Range       string `toml:"range" json:"range"`
}

type CSVConfig struct {
	Header           bool   `toml:"header" json:"header"`
	Separator        string `toml:"separator" json:"separator"`
	Terminator       string `toml:"terminator" json:"terminator"`
	Delimiter        string `toml:"delimiter" json:"delimiter"`
	EscapeBackslash  bool   `toml:"escape-backslash" json:"escape-backslash"`
	Charset          string `toml:"charset" json:"charset"`
	Rows             int    `toml:"rows" json:"rows"`
	OutputDir        string `toml:"output-dir" json:"output-dir"`
	TaskThreads      int    `toml:"task-threads" json:"task-threads"`
	TableThreads     int    `toml:"table-threads" json:"table-threads"`
	SQLThreads       int    `toml:"sql-threads" json:"sql-threads"`
	EnableCheckpoint bool   `toml:"enable-checkpoint" json:"enable-checkpoint"`
}

type FullConfig struct {
	ChunkSize        int  `form:"chunk-size" toml:"chunk-size" json:"chunk-size"`
	TaskThreads      int  `form:"task-threads" toml:"task-threads" json:"task-threads"`
	TableThreads     int  `form:"table-threads" toml:"table-threads" json:"table-threads"`
	SQLThreads       int  `form:"sql-threads" toml:"sql-threads" json:"sql-threads"`
	ApplyThreads     int  `form:"apply-threads" toml:"apply-threads" json:"apply-threads"`
	EnableCheckpoint bool `form:"enable-checkpoint" toml:"enable-checkpoint" json:"enable-checkpoint"`
}

type AllConfig struct {
	LogminerQueryTimeout int `toml:"logminer-query-timeout" json:"logminer-query-timeout"`
	FilterThreads        int `toml:"filter-threads" json:"filter-threads"`
	ApplyThreads         int `toml:"apply-threads" json:"apply-threads"`
	WorkerQueue          int `toml:"worker-queue" json:"worker-queue"`
	WorkerThreads        int `toml:"worker-threads" json:"worker-threads"`
}

type SourceConfig struct {
	Username      string   `form:"username" toml:"username" json:"username"`
	Password      string   `form:"password" toml:"password" json:"password"`
	Host          string   `form:"host" toml:"host" json:"host"`
	Port          int      `form:"port" toml:"port" json:"port"`
	ServiceName   string   `form:"service-name" toml:"service-name" json:"service-name"`
	ConnectParams string   `form:"connect-params" toml:"connect-params" json:"connect-params"`
	SessionParams []string `form:"session-params" toml:"session-params" json:"session-params"`
	SchemaName    string   `form:"schema-name" toml:"schema-name" json:"schema-name"`
	IncludeTable  []string `form:"include-table" toml:"include-table" json:"include-table"`
	ExcludeTable  []string `form:"exclude-table" toml:"exclude-table" json:"exclude-table"`
}

//type SourceConfig struct {
//	Username      string   `toml:"username" json:"username"`
//	Password      string   `toml:"password" json:"password"`
//	Host          string   `toml:"host" json:"host"`
//	Port          int      `toml:"port" json:"port"`
//	ServiceName   string   `toml:"service-name" json:"service-name"`
//	ConnectParams string   `toml:"connect-params" json:"connect-params"`
//	SessionParams []string `toml:"session-params" json:"session-params"`
//	SchemaName    string   `toml:"schema-name" json:"schema-name"`
//	IncludeTable  []string `toml:"include-table" json:"include-table"`
//	ExcludeTable  []string `toml:"exclude-table" json:"exclude-table"`
//}

//type TargetConfig struct {
//	Username      string `toml:"username" json:"username"`
//	Password      string `toml:"password" json:"password"`
//	Host          string `toml:"host" json:"host"`
//	Port          int    `toml:"port" json:"port"`
//	ConnectParams string `toml:"connect-params" json:"connect-params"`
//	MetaSchema    string `toml:"meta-schema" json:"meta-schema"`
//	SchemaName    string `toml:"schema-name",json:"schema-name"`
//	Overwrite     bool   `toml:"overwrite" json:"overwrite"`
//}

type TargetConfig struct {
	Username      string `form:"username" toml:"username" json:"username" binding:"required"`
	Password      string `form:"password" toml:"password" json:"password" binding:"required"`
	Host          string `form:"host" toml:"host" json:"host" binding:"required"`
	Port          int    `form:"port" toml:"port" json:"port" binding:"required"`
	ConnectParams string `form:"connect-params" toml:"connect-params" toml:"connect-params" json:"connect-params" binding:"required"`
	MetaSchema    string `form:"meta-schema" toml:"meta-schema" json:"meta-schema"`
	SchemaName    string `form:"schema-name" toml:"schema-name" json:"schema-name" binding:"required"`
	Overwrite     bool   `form:"overwrite" toml:"overwrite" json:"overwrite"`
}

type LogConfig struct {
	LogLevel   string `toml:"log-level" json:"log-level"`
	LogFile    string `toml:"log-file" json:"log-file"`
	MaxSize    int    `toml:"max-size" json:"max-size"`
	MaxDays    int    `toml:"max-days" json:"max-days"`
	MaxBackups int    `toml:"max-backups" json:"max-backups"`
}

// 读取配置文件
func ReadConfigFile(file string) (*CfgFile, error) {
	cfg := &CfgFile{}
	if err := cfg.configFromFile(file); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// 加载配置文件并解析
func (c *CfgFile) configFromFile(file string) error {
	if _, err := toml.DecodeFile(file, c); err != nil {
		return fmt.Errorf("failed decode toml config file %s: %v", file, err)
	}
	return nil
}

// 根据配置文件获取表列表
func (c *CfgFile) GenerateTables(engine *Engine) ([]string, error) {
	startTime := time.Now()
	var (
		exporterTableSlice []string
		err                error
	)
	switch {
	case len(c.SourceConfig.IncludeTable) != 0 && len(c.SourceConfig.ExcludeTable) == 0:
		if err := engine.IsExistOracleTable(c.SourceConfig.SchemaName, c.SourceConfig.IncludeTable); err != nil {
			return exporterTableSlice, err
		}
		exporterTableSlice = append(exporterTableSlice, c.SourceConfig.IncludeTable...)
	case len(c.SourceConfig.IncludeTable) == 0 && len(c.SourceConfig.ExcludeTable) != 0:
		exporterTableSlice, err = engine.FilterDifferenceOracleTable(c.SourceConfig.SchemaName, c.SourceConfig.ExcludeTable)
		if err != nil {
			return exporterTableSlice, err
		}
	case len(c.SourceConfig.IncludeTable) == 0 && len(c.SourceConfig.ExcludeTable) == 0:
		exporterTableSlice, err = engine.GetOracleTable(c.SourceConfig.SchemaName)
		if err != nil {
			return exporterTableSlice, err
		}
	default:
		return exporterTableSlice, fmt.Errorf("source config params include-table/exclude-table cannot exist at the same time")
	}

	if len(exporterTableSlice) == 0 {
		return exporterTableSlice, fmt.Errorf("exporter table slice can not null from reverse task")
	}
	endTime := time.Now()
	Logger.Info("get oracle to mysql all tables",
		zap.String("schema", c.SourceConfig.SchemaName),
		zap.Strings("tables", exporterTableSlice),
		zap.String("cost", endTime.Sub(startTime).String()))
	return exporterTableSlice, nil
}

func (c *CfgFile) String() string {
	cfg, err := json.Marshal(c)
	if err != nil {
		return "<nil>"
	}
	return string(cfg)
}
