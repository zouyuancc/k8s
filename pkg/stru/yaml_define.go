package yaml_define

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
