package dp

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	getcfg "k8s/pkg/getconfig"
	yaml_define "k8s/pkg/stru"
)

func Operate(data *yaml_define.Yaml) {
	Create(data)
	//Update(data,config)
}

//创建deployment
func Create(data *yaml_define.Yaml) {
	config := getcfg.Getconfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
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
									ContainerPort: data.Spec.Template.Spec.Containers[0].Ports.Containerport,
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
	fmt.Printf("Create deployment %q.\n", result.GetObjectMeta().GetName())
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
