package common

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s/cores"
	"strconv"
)

//创建deployment
func CreateDeployment(data *cores.Yaml) {
	clientset := cores.Getset()
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	fmt.Println("creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), dp_trans_to_kubernetes_struct(data), metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Deployment %q created. \n", result.GetObjectMeta().GetName()) //%q 单引号围绕的字符字面值，由Go语法安全地转义
}

//更新deployment
func UpdateDeployment(data *cores.Yaml) {
	clientset := cores.Getset()
	fmt.Println("updating deployment...")
	namespace := data.Metadata.Namespace
	deployment := dp_trans_to_kubernetes_struct(data)
	deployment, err := clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Deployment %q updated. \n", deployment.Name) //%q 单引号围绕的字符字面值，由Go语法安全地转义
}

//删除deployment
func DeleteDeployment(data *cores.Yaml) {
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("deleting deployment...")
}

//判断操作的deployment是否存在
func DeploymentExistJudge(data *cores.Yaml) bool {
	clientset := cores.Getset()
	deploymentName := data.Metadata.Name
	_, err := clientset.AppsV1().Deployments(data.Metadata.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func dp_trans_to_kubernetes_struct(data *cores.Yaml) *appsv1.Deployment {
	revcontainer := []apiv1.Container{}
	for i, v := range data.Spec.Template.Spec.Containers {
		revConPort := []apiv1.ContainerPort{}
		for _, p := range data.Spec.Template.Spec.Containers[i].Ports {
			tmpport := apiv1.ContainerPort{
				Name:          data.Metadata.Name + "-c" + strconv.Itoa(i) + "-p" + strconv.Itoa(i),
				HostPort:      p.HostPort,
				ContainerPort: p.ContainerPort,
				Protocol:      apiv1.Protocol(p.Protocol),
				HostIP:        p.HostIP,
			}
			revConPort = append(revConPort, tmpport)
		}
		tempcon := apiv1.Container{
			Name:       data.Metadata.Name + "-c" + strconv.Itoa(i),
			Image:      v.Image,
			Command:    v.Command,
			Args:       v.Args,
			WorkingDir: v.WorkingDir,
			Ports:      revConPort,
		}
		revcontainer = append(revcontainer, tempcon)
	}

	var dpobj metav1.ObjectMeta = metav1.ObjectMeta{
		Name:      data.Metadata.Name,
		Namespace: data.Metadata.Namespace,
		Labels:    data.Metadata.Labels,
	}
	var tmpobj metav1.ObjectMeta = metav1.ObjectMeta{
		Name:   data.Metadata.Name + "-c0",
		Labels: data.Spec.Template.Metadata.Labels,
	}

	var tmpspec apiv1.PodSpec = apiv1.PodSpec{
		Containers: revcontainer,
	}

	var dpspec appsv1.DeploymentSpec = appsv1.DeploymentSpec{
		Replicas: &data.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: data.Spec.Selector,
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: tmpobj,
			Spec:       tmpspec,
		},
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: dpobj,
		Spec:       dpspec,
	}
	return deployment
}
