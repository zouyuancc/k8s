package main

import (
	"context"
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	yaml "gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"log"
	"path/filepath"
	"strconv"
)

func main() {
	done := make(chan bool)
	go startServer(20000, done)
	<-done
}

//启动zmqserver端
func startServer(port int, done chan bool) {
	// REP表示server端
	socket, _ := zmq.NewSocket(zmq.REP)
	socket.Bind("tcp://127.0.0.1:" + strconv.Itoa(port))
	defer socket.Close()
	for {
		//Recv 和 Send必须交替进行
		resp, _ := socket.Recv(0)
		go parseargs([]byte(resp))
		socket.Send("Hello "+resp, 0)
	}
	done <- true
}

//解析从client端收过来的信息
func parseargs(resp []byte) {
	conf := new(Yaml)
	yaml.Unmarshal(resp, conf)
	//fmt.Println(conf)
	//createSource(conf)
	data, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data:\t", string(data))
}

//利用解析得到的结构体信息，创建deployment demo
func createSource(data *Yaml) {
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
	tmpport := data.Spec.Template.Spec.Containers[0].Ports.ContainerPort
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
						ContainerPort int32
						Other         int32
						Other1        int32
					}
				}
			}
		}
	}
}
