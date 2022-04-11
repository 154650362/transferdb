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
package main

import (
	"fmt"
	"github.com/wentaojin/transferdb/routers"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/pkg/errors"
	"github.com/wentaojin/transferdb/service"

	"github.com/wentaojin/transferdb/pkg/signal"

	"go.uber.org/zap"

	. "github.com/wentaojin/transferdb/conf"
)

//var (
//conf = flag.String("config", "config.toml", "specify the configuration file, default is config.toml")
//mode    = flag.String("mode", "", "specify the program running mode: [prepare reverse gather full csv all check diff]")
//version = flag.Bool("version", false, "view transferdb version info")
//)

func main() {
	//flag.Parse()

	// 获取程序版本
	//service.GetAppVersion(*version)

	// 读取配置文件
	//cfg, err := service.ReadConfigFile(*conf)
	//if err != nil {
	//	log.Fatalf("read config file [%s] failed: %v", *conf, err)
	//}

	go func() {
		if err := http.ListenAndServe(Gcfg.AppConfig.PprofPort, nil); err != nil {
			service.Logger.Fatal("listen and serve pprof failed", zap.Error(errors.Cause(err)))
		}
		os.Exit(0)
	}()

	// 初始化日志 logger
	if err := service.NewZapLogger(Gcfg); err != nil {
		log.Fatalf("create global zap logger failed: %v", err)
	}
	service.RecordAppVersion("transferdb", service.Logger, Gcfg)

	// 信号量监听处理
	signal.SetupSignalHandler(func() {
		os.Exit(1)
	})

	// 程序运行
	//if err = server.Run(cfg, *mode); err != nil {
	//	service.Logger.Fatal("server run failed", zap.Error(errors.Cause(err)))
	//}
	// 启动web server
	//todo 启动注册自己的的信息：addr，和当前时间
	// 在 service目录下的heartbeat实现
	router := routers.InitRouter()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", Gcfg.AppConfig.HTTPPort),
		Handler: router,
		//ReadTimeout:    setting.ReadTimeout,
		//WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		service.Logger.Fatal("http server run failed", zap.Error(errors.Cause(err)))
	}
}
