package conf

import (
	"github.com/wentaojin/transferdb/service"
	"log"
	"sync"
)

// 初始化配置文件
var Gcfg *service.CfgFile
var defaultConfig = "conf/configcopy.toml"
var once sync.Once

func init() {
	once.Do(initconf)
}

func initconf() {
	var err error
	Gcfg, err = service.ReadConfigFile(defaultConfig)
	if err != nil {
		log.Fatalf("read config file [%s] failed: %v", defaultConfig, err)
	}
	log.Printf("Gcfg is %v", Gcfg)
}

func Getconf() {

}
