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
package db

import (
	"fmt"
	"strings"

	"github.com/WentaoJin/dbsyncer/zlog"
	"go.uber.org/zap"

	"github.com/WentaoJin/dbsyncer/util"
)

func (e *Engine) IsExistMySQLSchema(schemaName string) (bool, error) {
	schemas, err := e.getMySQLSchema()
	if err != nil {
		return false, err
	}
	if !util.IsContainString(schemas, strings.ToUpper(schemaName)) {
		return false, nil
	}
	return true, nil
}

func (e *Engine) FilterIntersectionMySQLTable(schemaName string, includeTables []string) ([]string, error) {
	tables, err := e.getMySQLTable(schemaName)
	if err != nil {
		return []string{}, err
	}
	var includeTbl []string
	for _, tbl := range includeTables {
		includeTbl = append(includeTbl, strings.ToUpper(tbl))
	}
	return util.FilterIntersectionStringItems(tables, includeTbl), nil
}

func (e *Engine) RenameMySQLTableName(schemaName string, tableName string) error {
	backupTable := fmt.Sprintf("%s_bak", tableName)
	querySQL := fmt.Sprintf("RENAME TABLE `%s`.`%s` TO `%s`.`%s`", schemaName, tableName, schemaName, backupTable)
	zlog.Logger.Info("Exec SQL",
		zap.String("schema", schemaName),
		zap.String("table", tableName),
		zap.String("sql", fmt.Sprintf("%v", querySQL)))
	_, _, err := Query(e.MysqlDB, querySQL)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) IsExistMysqlIndex(schemaName, tableName, indexName string) bool {
	querySQL := fmt.Sprintf(`SELECT count(1) AS CT
FROM information_schema.statistics 
WHERE upper(table_schema) = upper('%s')
AND upper(table_name) = upper('%s')
AND upper(index_name) = upper('%s')`, schemaName, tableName, indexName)
	_, res, _ := Query(e.MysqlDB, querySQL)
	if res[0]["CT"] == "0" {
		return false
	}
	return true
}

func (e *Engine) IsExistOracleSchema(schemaName string) error {
	schemas, err := e.getOracleSchema()
	if err != nil {
		return err
	}
	if !util.IsContainString(schemas, strings.ToUpper(schemaName)) {
		return fmt.Errorf("oracle schema [%s] isn't exist in the database", schemaName)
	}
	return nil
}

func (e *Engine) IsExistOracleTable(schemaName string, includeTables []string) error {
	tables, err := e.getOracleTable(schemaName)
	if err != nil {
		return err
	}
	ok, noExistTables := util.IsSubsetString(tables, includeTables)
	if !ok {
		return fmt.Errorf("oracle include-tables values [%v] isn't exist in the db schema [%v]", noExistTables, schemaName)
	}
	return nil
}

func (e *Engine) FilterDifferenceOracleTable(schemaName string, excludeTables []string) ([]string, error) {
	tables, err := e.getOracleTable(schemaName)
	if err != nil {
		return []string{}, err
	}
	return util.FilterDifferenceStringItems(tables, excludeTables), nil
}

func (e *Engine) GetOracleTableComment(schemaName string, tableName string) ([]map[string]string, error) {
	var (
		comments []map[string]string
		err      error
	)
	querySQL := fmt.Sprintf(`select table_name,table_type,comments 
from all_tab_comments 
where 
table_type = 'TABLE'
and upper(owner)=upper('%s')
and upper(table_name)=upper('%s')`, strings.ToUpper(schemaName), strings.ToUpper(tableName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return comments, err
	}
	return res, nil
}

func (e *Engine) GetOracleTableColumn(schemaName string, tableName string) ([]map[string]string, error) {
	querySQL := fmt.Sprintf(`select t.COLUMN_NAME,
	     t.DATA_TYPE,
	     NVL(t.DATA_LENGTH,0) AS DATA_LENGTH,
	     NVL(t.DATA_PRECISION,0) AS DATA_PRECISION,
	     NVL(t.DATA_SCALE,0) AS DATA_SCALE,
	     t.NULLABLE,
	     t.DATA_DEFAULT,
	     c.COMMENTS
	from all_tab_columns t, all_col_comments c
	where t.table_name = c.table_name
	 and t.column_name = c.column_name
     and t.owner = c.owner
	 and t.table_name = upper('%s')
	 and t.owner = upper('%s')`,
		strings.ToUpper(tableName),
		strings.ToUpper(schemaName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (e *Engine) GetOracleTablePrimaryKey(schemaName string, tableName string) ([]map[string]string, error) {
	// for the primary key of an Engine table, you can use the following command to set whether the primary key takes effect.
	// disable the primary key: alter table tableName disable primary key;
	// enable the primary key: alter table tableName enable primary key;
	// primary key status Disabled will not do primary key processing
	querySQL := fmt.Sprintf(`select cu.constraint_name,
       LISTAGG(cu.column_name, ',') WITHIN GROUP(ORDER BY cu.POSITION) AS column_list
  from all_cons_columns cu, all_constraints au
 where cu.constraint_name = au.constraint_name
   and au.constraint_type = 'P'
   and au.STATUS = 'ENABLED'
   and upper(au.table_name) = upper('%s')
   and upper(cu.owner) = upper('%s')
 group by cu.constraint_name`,
		strings.ToUpper(tableName),
		strings.ToUpper(schemaName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (e *Engine) GetOracleTableUniqueKey(schemaName string, tableName string) ([]map[string]string, error) {
	querySQL := fmt.Sprintf(`select cu.constraint_name,
       LISTAGG(cu.column_name, ',') WITHIN GROUP(ORDER BY cu.POSITION) AS column_list
  from all_cons_columns cu, all_constraints au
 where cu.constraint_name = au.constraint_name
   and au.constraint_type = 'U'
   and au.STATUS = 'ENABLED'
   and upper(au.table_name) = upper('%s')
   and upper(cu.owner) = upper('%s')
 group by cu.constraint_name`,
		strings.ToUpper(tableName),
		strings.ToUpper(schemaName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (e *Engine) GetOracleTableForeignKey(schemaName string, tableName string) ([]map[string]string, error) {
	querySQL := fmt.Sprintf(`select t1.table_name,
       t2.table_name        as RTABLE_NAME,
       t1.constraint_name,
       -- t1.r_constraint_name as RCONSTRAINT_NAME,
       a1.column_name,
       a2.column_name       as RCOLUMN_NAME
  from all_constraints  t1,
       all_constraints  t2,
       all_cons_columns a1,
       all_cons_columns a2
 where t1.r_constraint_name = t2.constraint_name
   and t1.constraint_name = a1.constraint_name
   and t1.r_constraint_name = a2.constraint_name
  and upper(t1.table_name) = upper('%s')
   and upper(t1.owner) = upper('%s')`,
		strings.ToUpper(tableName),
		strings.ToUpper(schemaName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (e *Engine) GetOracleTableIndex(schemaName string, tableName string) ([]map[string]string, error) {
	querySQL := fmt.Sprintf(`select T.TABLE_NAME,
       I.UNIQUENESS, --是否唯一索引
       T.INDEX_NAME,
       --T.COLUMN_POSITION,
       LISTAGG(T.COLUMN_NAME, ',') WITHIN GROUP(ORDER BY T.COLUMN_POSITION) AS column_list
  FROM ALL_IND_COLUMNS T, ALL_INDEXES I, ALL_CONSTRAINTS C
 WHERE T.INDEX_NAME = I.INDEX_NAME
   AND T.INDEX_NAME = C.CONSTRAINT_NAME(+)
   AND I.INDEX_TYPE != 'FUNCTION-BASED NORMAL' --排除基于函数的索引
   AND I.INDEX_TYPE != 'BITMAP' --排除位图索引
   AND C.CONSTRAINT_TYPE is Null --排除主键、唯一约束索引
   AND T.TABLE_NAME = upper('%s')
   AND T.TABLE_OWNER = upper('%s')
 group by T.TABLE_NAME,
          I.UNIQUENESS, --是否唯一索引
          T.INDEX_NAME`,
		strings.ToUpper(tableName),
		strings.ToUpper(schemaName))
	_, res, err := Query(e.OracleDB, querySQL)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (e *Engine) getOracleSchema() ([]string, error) {
	var (
		schemas []string
		err     error
	)
	cols, res, err := Query(e.OracleDB, `SELECT DISTINCT username FROM ALL_USERS`)
	if err != nil {
		return schemas, err
	}
	for _, col := range cols {
		for _, r := range res {
			schemas = append(schemas, strings.ToUpper(r[col]))
		}
	}
	return schemas, nil
}

func (e *Engine) getOracleTable(schemaName string) ([]string, error) {
	var (
		tables []string
		err    error
	)
	cols, res, err := Query(e.OracleDB, fmt.Sprintf(`SELECT table_name FROM ALL_TABLES WHERE UPPER(owner) = UPPER('%s')`, schemaName))
	if err != nil {
		return tables, err
	}
	for _, col := range cols {
		for _, r := range res {
			tables = append(tables, r[col])
		}
	}
	return tables, nil
}

func (e *Engine) getMySQLSchema() ([]string, error) {
	var (
		schemas []string
		err     error
	)
	cols, res, err := Query(e.MysqlDB, `SELECT DISTINCT(schema_name) AS SCHEMA_NAME FROM information_schema.SCHEMATA`)
	if err != nil {
		return schemas, err
	}
	for _, col := range cols {
		for _, r := range res {
			schemas = append(schemas, strings.ToUpper(r[col]))
		}
	}
	return schemas, nil
}

func (e *Engine) getMySQLTable(schemaName string) ([]string, error) {
	var (
		tables []string
		err    error
	)
	cols, res, err := Query(e.MysqlDB, fmt.Sprintf(`select table_name from information_schema.tables where upper(table_schema) = upper('%s') and upper(table_type)=upper('base table')`, schemaName))
	if err != nil {
		return tables, err
	}
	for _, col := range cols {
		for _, r := range res {
			tables = append(tables, strings.ToUpper(r[col]))
		}
	}
	return tables, nil
}
