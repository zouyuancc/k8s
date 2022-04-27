package dp

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s/pkg/getclientset"
	yaml_define "k8s/pkg/stru"
)

//创建deployment
func Create(data *yaml_define.Yaml) {
	clientset := getclientset.Getset()
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Metadata.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Spec.Selector,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Spec.Template.Metadata.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  data.Spec.Template.Spec.Containers[0].Name,
							Image: data.Spec.Template.Spec.Containers[0].Image,
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: data.Spec.Template.Spec.Containers[0].Ports.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Deployment %q created. \n", result.GetObjectMeta().GetName()) //%q 单引号围绕的字符字面值，由Go语法安全地转义
}

//更新deployment
func Update(data *yaml_define.Yaml) {
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("updating deployment...")
}

//删除deployment
func Delete(data *yaml_define.Yaml) {
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("deleting deployment...")
}

//判断操作的deployment是否存在
func Existjudge(data *yaml_define.Yaml) bool {
	clientset := getclientset.Getset()
	deploymentName := data.Metadata.Name
	_, err := clientset.AppsV1().Deployments(data.Metadata.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}
