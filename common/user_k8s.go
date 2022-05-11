package common

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"gopkg.in/yaml.v2"
	"k8s/cores"
	"time"
)

type User struct {
	Name   string
	Socket *zmq.Socket
	Server *Server
}

type Job struct {
	Time     time.Time
	Kind     string
	TaskName string
}

type History struct {
	User string
	Task []*Job
}

func NewUser(socket *zmq.Socket, server *Server) *User {
	user := &User{
		Server: server,
		Socket: socket,
	}

	return user
}

func (user *User) Online() {
	if _, ok := user.Server.OnlineMap[user.Name]; ok {
		newJob := &Job{
			Time: time.Now(),
		}
		user.Server.Maplock.Lock()
		user.Server.OnlineMap[user.Name].Task = append(user.Server.OnlineMap[user.Name].Task, newJob)
		user.Server.Maplock.Unlock()
		return
	}
	user.Server.Maplock.Lock()
	newJob := &Job{
		Time: time.Now(),
	}
	task := []*Job{newJob}
	history := &History{
		Task: task,
	}
	user.Server.OnlineMap[user.Name] = history
	user.Server.Maplock.Unlock()
}

//填充初始化
func Fillstructure(user *User, data *cores.Yaml) {
	user.Name = data.User
	user.Online()
	user.Server.OnlineMap[user.Name].User = user.Name
	task := user.Server.OnlineMap[user.Name].Task
	task[len(task)-1].Kind = data.Kind
	task[len(task)-1].TaskName = data.Metadata.Name
}

//解析从client端收过来的信息
func (user *User) Parseargs(resp []byte) {
	strback := ""
	data := new(cores.Yaml)
	yaml.Unmarshal(resp, data)

	//用户上线
	Fillstructure(user, data)

	switch data.Operation {
	case "apply":
		switch data.Kind {
		case "Deployment":
			if DeploymentExistJudge(data) {
				UpdateDeployment(data)
			} else {
				CreateDeployment(data)
			}
			break
		case "Service":
			if ServiceExistJudge(data) {
				UpdateService(data)
			} else {
				CreateService(data)
			}
			break
		}
		break
	case "delete":
		switch data.Kind {
		case "Deployment":
			if DeploymentExistJudge(data) {
				DeleteDeployment(data)
			} else {
				fmt.Printf("Deployment %q not exists\n", data.Metadata.Name)
			}
			break
		case "Service":
			if ServiceExistJudge(data) {
				DeleteService(data)
			} else {
				fmt.Printf("Service %q not exists\n", data.Metadata.Name)
			}
			break
		}
		break
	case "list":
		switch data.Kind {
		case "deployment":
			strback += DeploymentList(data)
			break
		case "service":
			strback += ServiceList(data)
			break
		}
		break

	}
	user.Socket.Send(strback, 0)
}
