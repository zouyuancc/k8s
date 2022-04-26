package parse

import (
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	dp "k8s/pkg/deployment"
	yaml_define "k8s/pkg/stru"
	"path/filepath"
)

//利用解析得到的结构体信息，判断资源类型
func CreateSource(data *yaml_define.Yaml) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatal(err)
		return
	}

	if data.Kind == "Deployment" {
		dp.Create(data, config)
	}
	if data.Kind == "Service" {
		dp.Operate(data, config)
	}
}
