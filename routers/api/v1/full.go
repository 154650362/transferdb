package v1

//case "full":
//// 全量数据 ETL 非一致性（基于某个时间点，而是直接基于现有 SCN）抽取，离线环境提供与原库一致性
//engine, err := NewEngineDB(
//cfg.SourceConfig, cfg.TargetConfig, cfg.AppConfig.SlowlogThreshold,
//cfg.FullConfig.TableThreads*cfg.FullConfig.SQLThreads*cfg.FullConfig.ApplyThreads)
//if err != nil {
//return err
//}
//if err = taskflow.FullSyncOracleTableRecordToMySQL(cfg, engine); err != nil {
//return err
//}

import "github.com/gin-gonic/gin"

//todo
func Full(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
