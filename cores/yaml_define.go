package cores

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//自定义结构体，用于解析client发过来的yaml文件
type Yaml struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name      string            `yaml:"name"`
		Namespace string            `yaml:"namespace"`
		Labels    map[string]string `yaml:"labels"`
	}
	Spec struct {
		Replicas int32 `yaml:"replicas"`
		Selector map[string]string
		Template struct {
			Metadata struct {
				Labels map[string]string
			}
			Spec struct {
				Containers []struct {
					Image string `yaml:"image"`
					Name  string `yaml:"name"`
					Ports []struct {
						Name          string
						HostPort      int32
						ContainerPort int32 `yaml:"containerPort"`
						Protocol      string
						HostIP        string
					}
					Command    []string `yaml:"command"`
					Args       []string
					WorkingDir string
					Resources  struct {
						Requests struct {
							Memory    string `yaml:"memory"`
							Cpu       string `yaml:"cpu"`
							NvidiaGpu string `yaml:"nvidia.com/gpu"`
						}
						Limits struct {
							Memory    string `yaml:"memory"`
							Cpu       string `yaml:"cpu"`
							NvidiaGpu string `yaml:"nvidia.com/gpu"`
						}
					}
					VolumeMount struct {
						Name      string `yaml:"name"`
						MountPath string `yaml:"mountPath"`
					}
				}
				Volumes struct {
					Name string `yaml:name`
					Nfs  struct {
						Path   string `yaml:path`
						Server string `yaml:server`
					}
				}
			}
		}

		//Service
		Type string `yaml:"type"`
		//Ports []v1.ServicePort

		Ports []struct {
			TargetPort intstr.IntOrString `yaml:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`
			Name       string             `yaml:"name"`
			Port       int32              `yaml:"port"`
			NodePort   int32              `yaml:"nodePort,omitempty" protobuf:"varint,5,opt,name=nodePort"`
			Protocol   v1.Protocol        `yaml:"protocol"`
		}
	}
	User      string `yaml:"user"`
	Operation string `yaml:"operation"`
}
