package main

//import (
//	"fmt"
//	appsv1 "k8s.io/api/apps/v1"
//	apiv1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/api/resource"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s/cores"
//	"strconv"
//)
//
//package common
//
//import (
//"context"
//"fmt"
//appsv1 "k8s.io/api/apps/v1"
//apiv1 "k8s.io/api/core/v1"
//"k8s.io/apimachinery/pkg/api/resource"
//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//"k8s/cores"
//"strconv"
//)
//
//type ResourceName string
//
////创建deployment
//func CreateDeployment(data *cores.Yaml) {
//	clientset := cores.Getset()
//	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
//	fmt.Println("creating deployment...")
//	result, err := deploymentsClient.Create(context.TODO(), dp_trans_to_kubernetes_struct(data), metav1.CreateOptions{})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("Deployment %q created. \n", result.GetObjectMeta().GetName()) //%q 单引号围绕的字符字面值，由Go语法安全地转义
//}
//
////更新deployment
//func UpdateDeployment(data *cores.Yaml) {
//	clientset := cores.Getset()
//	fmt.Printf("updating deployment %q\n", data.Metadata.Name)
//	namespace := data.Metadata.Namespace
//	deployment := dp_trans_to_kubernetes_struct(data)
//	deployment, err := clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("Deployment %q updated. \n", deployment.Name) //%q 单引号围绕的字符字面值，由Go语法安全地转义
//}
//
////删除deployment
//func DeleteDeployment(data *cores.Yaml) {
//	clientset := cores.Getset()
//	fmt.Printf("deleting deployment %q\n", data.Metadata.Name)
//	if err := clientset.AppsV1().Deployments(data.Metadata.Namespace).Delete(context.TODO(), data.Metadata.Name, metav1.DeleteOptions{}); err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("Deleted deployment %q\n", data.Metadata.Name)
//}
//
////判断操作的deployment是否存在
//func DeploymentExistJudge(data *cores.Yaml) bool {
//	clientset := cores.Getset()
//	deploymentName := data.Metadata.Name
//	_, err := clientset.AppsV1().Deployments(data.Metadata.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
//	if err != nil {
//		return false
//	}
//	return true
//}
//
////查看deployment列表
//func DeploymentList(data *cores.Yaml) string {
//	dplist := ""
//	clientset := cores.Getset()
//	namespace := data.Metadata.Namespace
//	deployments, _ := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
//	dplist += "deploymnet_name\t\t\tREADY\t\tUP-TO-DATE\t\tAVAILABLE\t\tAGE\n"
//	for _, deployment := range deployments.Items {
//		dplist += fmt.Sprintf("%-9.8s\t\t\t%d/%d\t\t%d\t\t\t%d\t\t\t%s\n",
//			//fmt.Printf("%-9.8s\t\t\t%d/%d\t\t%d\t\t\t%d\t\t\t%s\n",
//			deployment.Name,
//			deployment.Status.ReadyReplicas, deployment.Status.Replicas,
//			deployment.Status.UpdatedReplicas,
//			deployment.Status.AvailableReplicas,
//			deployment.CreationTimestamp,
//		)
//	}
//	return dplist
//}
//
//func dp_trans_to_kubernetes_struct(data *cores.Yaml) *appsv1.Deployment {
//	revcontainer := []apiv1.Container{}
//	for i, v := range data.Spec.Template.Spec.Containers {
//		revConPort := []apiv1.ContainerPort{}
//		for _, p := range data.Spec.Template.Spec.Containers[i].Ports {
//			tmpport := apiv1.ContainerPort{
//				Name:          data.Metadata.Name + "-c" + strconv.Itoa(i) + "-p" + strconv.Itoa(i),
//				HostPort:      p.HostPort,
//				ContainerPort: p.ContainerPort,
//				Protocol:      apiv1.Protocol(p.Protocol),
//				HostIP:        p.HostIP,
//			}
//			revConPort = append(revConPort, tmpport)
//		}
//
//		//resource
//		var res apiv1.ResourceRequirements
//		resourceLimit := make(map[apiv1.ResourceName]resource.Quantity)
//		var tmpcpu apiv1.ResourceName = "cpu"
//		var tmpmem apiv1.ResourceName = "memory"
//		var tmpgpu apiv1.ResourceName = "nvidia.com/gpu"
//
//		resourceLimit[tmpcpu] = resource.MustParse(string(data.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu))
//		resourceLimit[tmpmem] = resource.MustParse(string(data.Spec.Template.Spec.Containers[i].Resources.Limits.Memory))
//		resourceLimit[tmpgpu] = resource.MustParse(string(data.Spec.Template.Spec.Containers[i].Resources.Limits.NvidiaGpu))
//		res.Limits = resourceLimit
//
//		//request
//		resourceRequest := make(map[apiv1.ResourceName]resource.Quantity)
//		resourceRequest[tmpcpu] = resource.MustParse(string(data.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu))
//		resourceRequest[tmpmem] = resource.MustParse(string(data.Spec.Template.Spec.Containers[i].Resources.Limits.Memory))
//		resourceRequest[tmpgpu] = resource.MustParse(string(data.Spec.Template.Spec.Containers[i].Resources.Limits.NvidiaGpu))
//		res.Requests = resourceRequest
//
//		//volumenounts
//		//volumemounts := []apiv1.VolumeMount{}
//		//volumem := &apiv1.VolumeMount{
//		//	Name:      data.Spec.Template.Spec.Containers[i].VolumeMount.Name,
//		//	MountPath: data.Spec.Template.Spec.Containers[i].VolumeMount.MountPath,
//		//}
//		//volumemounts = append(volumemounts, *volumem)
//
//		tempcon := apiv1.Container{
//			Name:       data.Metadata.Name + "-c" + strconv.Itoa(i),
//			Image:      v.Image,
//			Command:    v.Command,
//			Args:       v.Args,
//			WorkingDir: v.WorkingDir,
//			Ports:      revConPort,
//			Resources:  res,
//			//VolumeMounts: volumemounts,
//		}
//		revcontainer = append(revcontainer, tempcon)
//	}
//
//	var tmpobj metav1.ObjectMeta = metav1.ObjectMeta{
//		Name:   data.Metadata.Name,
//		Labels: data.Spec.Template.Metadata.Labels,
//	}
//
//	//volumes
//	//volumes := []apiv1.Volume{}
//	//tmpnfs := &apiv1.NFSVolumeSource{
//	//	Path:   data.Spec.Template.Spec.Volumes.Nfs.Path,
//	//	Server: data.Spec.Template.Spec.Volumes.Nfs.Server,
//	//}
//	//volumeSource := apiv1.VolumeSource{
//	//	NFS: tmpnfs,
//	//}
//	//volume := apiv1.Volume{
//	//	Name:         data.Spec.Template.Spec.Volumes.Name,
//	//	VolumeSource: volumeSource,
//	//}
//	//volumes = append(volumes, volume)
//	//var tmpobj metav1.ObjectMeta = metav1.ObjectMeta{
//	//	Name:   data.Metadata.Name,
//	//	Labels: data.Spec.Template.Metadata.Labels,
//	//}
//
//	var tmpspec apiv1.PodSpec = apiv1.PodSpec{
//		Hostname:   data.Metadata.Namespace + "pod",
//		Containers: revcontainer,
//		//Volumes:    volumes,
//	}
//
//	var dpobj metav1.ObjectMeta = metav1.ObjectMeta{
//		Name:      data.Metadata.Name,
//		Namespace: data.Metadata.Namespace,
//		Labels:    data.Metadata.Labels,
//	}
//	var dpspec appsv1.DeploymentSpec = appsv1.DeploymentSpec{
//		Replicas: &data.Spec.Replicas,
//		Selector: &metav1.LabelSelector{
//			MatchLabels: data.Spec.Selector,
//		},
//		Template: apiv1.PodTemplateSpec{
//			ObjectMeta: tmpobj,
//			Spec:       tmpspec,
//		},
//	}
//	deployment := &appsv1.Deployment{
//		ObjectMeta: dpobj,
//		Spec:       dpspec,
//	}
//	return deployment
//}
//
