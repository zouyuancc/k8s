package yaml_define

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
					Ports struct {
						ContainerPort int32 `yaml:"containerPort"`
					}
				}
			}
		}
	}
}
