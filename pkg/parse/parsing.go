package parse

import (
	dp "k8s/pkg/deployment"
	yaml_define "k8s/pkg/stru"
)

//利用解析得到的结构体信息，判断资源类型
func OperateSource(data *yaml_define.Yaml) {
	//judge to decide to create or update
	if data.Kind == "Deployment" {
		if dp.Existjudge(data) {
			dp.Update(data)
		} else {
			dp.Create(data)
		}
	}
	if data.Kind == "Service" {
	}
}