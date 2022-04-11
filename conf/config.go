package conf

import "github.com/wentaojin/transferdb/service"
import "log"

// 初始化配置文件
var Gcfg *service.CfgFile
var defaultConfig = "configcopy.toml"

func Read2conf() {
	Gcfg, err := service.ReadConfigFile(defaultConfig)
	if err != nil {
		log.Fatalf("read config file [%s] failed: %v", defaultConfig, err)
	}
	log.Printf("Gcfg is %v", Gcfg)
}
