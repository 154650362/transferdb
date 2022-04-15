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

	"github.com/wentaojin/transferdb/service"

	"github.com/wentaojin/transferdb/server"
)

func main() {
	mysqlCfg := service.TargetConfig{
		Username:      "root",
		Password:      "123456",
		Host:          "192.168.1.112",
		Port:          3306,
		ConnectParams: "charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true&tidb_txn_mode='optimistic'",
		MetaSchema:    "todb",
	}
	engine, err := server.NewMySQLEngineGeneralDB(mysqlCfg, 300, 10)
	if err != nil {
		fmt.Println(err)
	}
	cols, res, err := service.Query(engine.MysqlDB, "select * from t1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cols)
	fmt.Println(res)
}
