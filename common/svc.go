package common

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s/cores"
)

func CreateService(data *cores.Yaml) {
	clientset := cores.Getset()
	fmt.Println("creating service...")
	service, err := clientset.CoreV1().Services(data.Metadata.Namespace).Create(context.TODO(), svc_trans_to_kubernetes_struct(data), metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Service %q created. \n", service.GetObjectMeta().GetName())
}

func UpdateService(data *cores.Yaml) {
	clientset := cores.Getset()
	fmt.Printf("Updating service %q\n", data.Metadata.Name)
	namespace := data.Metadata.Namespace
	service := svc_trans_to_kubernetes_struct(data)
	service, err := clientset.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Service %q updated.\n", service.Name)
}

func DeleteService(data *cores.Yaml) {

}

//判断操作的deployment是否存在
func ServiceExistJudge(data *cores.Yaml) bool {
	clientset := cores.Getset()
	serviceName := data.Metadata.Name
	_, err := clientset.CoreV1().Services(data.Metadata.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func svc_trans_to_kubernetes_struct(data *cores.Yaml) *v1.Service {
	specports := []v1.ServicePort{}
	for _, v := range data.Spec.Ports {
		port := v1.ServicePort{
			Name:       v.Name,
			Protocol:   v.Protocol,
			Port:       v.Port,
			TargetPort: v.TargetPort,
			NodePort:   v.NodePort,
		}
		specports = append(specports, port)
	}

	var svcobj metav1.ObjectMeta = metav1.ObjectMeta{
		Name:      data.Metadata.Name,
		Namespace: data.Metadata.Namespace,
		Labels:    data.Spec.Template.Metadata.Labels,
	}

	var svcspec v1.ServiceSpec = v1.ServiceSpec{
		Type:     v1.ServiceType(data.Spec.Type),
		Selector: data.Spec.Selector,
		Ports:    specports,
	}

	service := &v1.Service{
		ObjectMeta: svcobj,
		Spec:       svcspec,
	}
	return service
}
