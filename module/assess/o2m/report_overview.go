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
package o2m

import (
	"context"
	"encoding/json"
	"github.com/wentaojin/transferdb/common"
	"github.com/wentaojin/transferdb/database/meta"
	"github.com/wentaojin/transferdb/database/oracle"
)

type ReportOverview struct {
	ReportName        string `json:"report_name"`
	ReportUser        string `json:"report_user"`
	HostName          string `json:"host_name"`
	PlatformName      string `json:"platform_name"`
	DBName            string `json:"db_name"`
	GlobalDBName      string `json:"global_db_name"`
	ClusterDB         string `json:"cluster_db"`
	ClusterDBInstance string `json:"cluster_db_instance"`
	InstanceName      string `json:"instance_name"`
	InstanceNumber    string `json:"instance_number"`
	ThreadNumber      string `json:"thread_number"`
	BlockSize         string `json:"block_size"`
	TotalUsedSize     string `json:"total_used_size"`
	HostCPUS          string `json:"host_cpus"`
	HostMem           string `json:"host_mem"`
	CharacterSet      string `json:"character_set"`
}

func (ro *ReportOverview) String() string {
	jsonStr, _ := json.Marshal(ro)
	return string(jsonStr)
}

type ReportSummary struct {
	AssessType    string `json:"assess_type"`
	AssessName    string `json:"assess_name"`
	AssessTotal   int    `json:"assess_total"`
	Compatible    int    `json:"compatible"`
	Incompatible  int    `json:"incompatible"`
	Convertible   int    `json:"convertible"`
	InConvertible int    `json:"inconvertible"`
}

func (ro *ReportSummary) String() string {
	jsonStr, _ := json.Marshal(ro)
	return string(jsonStr)
}

/*
Oracle Database Overview
*/
func GetAssessDatabaseOverviewResult(ctx context.Context, metaDB *meta.Meta, oracle *oracle.Oracle, reportName, reportUser string) (*ReportOverview, ReportSummary, error) {
	// 获取自定义兼容性内容
	compatibles, err := meta.NewBuildinObjectCompatibleModel(metaDB).BatchQueryObjAssessCompatible(ctx, &meta.BuildinObjectCompatible{
		DBTypeS: common.TaskDBOracle,
		DBTypeT: common.TaskDBMySQL,
	})
	if err != nil {
		return nil, ReportSummary{}, err
	}
	objAssessCompsMap := make(map[string]meta.BuildinObjectCompatible)
	for _, c := range compatibles {
		objAssessCompsMap[common.StringUPPER(c.ObjectNameS)] = c
	}

	overview, rs, err := AssessOracleDBOverview(oracle, objAssessCompsMap, reportName, reportUser)
	if err != nil {
		return nil, ReportSummary{}, err
	}

	return overview, rs, nil
}
