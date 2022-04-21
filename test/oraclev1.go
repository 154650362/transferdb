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
	oraCfg := service.SourceConfig{
		Username:      "test",
		Password:      "test",
		Host:          "192.168.31.122",
		Port:          1521,
		ServiceName:   "helowin",
		ConnectParams: "poolMinSessions=10&poolMaxSessions=1000&poolWaitTimeout=60s&poolSessionMaxLifetime=1h&poolSessionTimeout=5m&poolIncrement=10&timezone=Local",
		SessionParams: []string{"ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD HH24:MI:SS'", "ALTER SESSION SET TIME_ZONE='Asia/Shanghai'"},
		SchemaName:    "test",
		IncludeTable:  nil,
		ExcludeTable:  nil,
	}
	sqlDB, err := server.NewOracleDBEngine(oraCfg)
	if err != nil {
		fmt.Println(err)
	}

	engine := service.Engine{
		OracleDB: sqlDB,
	}

	//_, _, err = service.Query(engine.OracleDB, "alter session set nls_date_format = 'yyyy-mm-dd hh24:mi:ss'")
	//if err != nil {
	//	fmt.Println(err)
	//}
	cols, res, err := service.Query(engine.OracleDB, "select * from test.marvin3 where rownum < 2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cols)

	fmt.Println(res)

	cols, res, err = service.Query(engine.OracleDB, "select to_char(sysdate,'yyyy-mm-dd HH24:mi:ss') from dual")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cols)

	fmt.Println(res)

}
