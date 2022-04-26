package parsedp

import (
	"context"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"path/filepath"
)

//利用解析得到的结构体信息，创建deployment demo
func CreateSource(data *Yaml) {
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
	//containerPort, _ := strconv.ParseInt(data.Spec.Template.Spec.Containers[0].Ports.ContainerPort, 10, 32)
	//tmpport := int32(containerPort)
	tmpport := data.Spec.Template.Spec.Containers[0].Ports.Containerport
	//tmpport := data.Spec.Template.Spec.Containers[0].Ports

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
									ContainerPort: tmpport,
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
		panic(err)
	}
	fmt.Printf("Create deployment %q.\n", result.GetObjectMeta().GetName())
}

//自定义结构体，用于解析client发过来的yaml文件
type Yaml struct {
	Kind     string
	Metadata struct {
		Name   string
		Labels map[string]string
	}
	Spec struct {
		Replicas int32
		Selector map[string]string
		Template struct {
			Metadata struct {
				Name   string
				Labels map[string]string
			}
			Spec struct {
				Containers []struct {
					Image string
					Name  string
					Ports struct {
						Containerport int32
					}
				}
			}
		}
	}
}
