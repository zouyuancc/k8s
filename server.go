package main

import (
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	yaml "gopkg.in/yaml.v2"
	"k8s/common"
	"k8s/cores"
	"strconv"
)

var ip string
var port int

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Input the server ip,default 127.0.0.1")
	flag.IntVar(&port, "port", 20000, "Input the server port,default 20000")
}

type Server struct {
	Ip   string
	Port int

	socket *zmq.Socket
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	socket, _ := zmq.NewSocket(zmq.REP)
	server.socket = socket
	server.socket.Bind("tcp://" + ip + ":" + strconv.Itoa(port))
	return server
}

func (server *Server) Start(done chan bool) {
	defer server.socket.Close()
	for {
		//Recv 和 Send必须交替进行
		resp, _ := server.socket.Recv(0)
		if resp != "" {
			go server.Parseargs([]byte(resp))
		}
	}
	done <- true

}

func main() {
	done := make(chan bool)
	server := NewServer(ip, port)
	go server.Start(done)
	//go startServer(20000, done)
	<-done
}

//解析从client端收过来的信息
func (server *Server) Parseargs(resp []byte) {
	strback := ""
	data := new(cores.Yaml)
	yaml.Unmarshal(resp, data)
	//lins, err := json.Marshal(data)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("data:\t", string(lins))

	switch data.Operation {
	case "apply":
		switch data.Kind {
		case "Deployment":
			if common.DeploymentExistJudge(data) {
				common.UpdateDeployment(data)
			} else {
				common.CreateDeployment(data)
			}
			break
		case "Service":
			if common.ServiceExistJudge(data) {
				common.UpdateService(data)
			} else {
				common.CreateService(data)
			}
			break
		}
		break
	case "delete":
		switch data.Kind {
		case "Deployment":
			if common.DeploymentExistJudge(data) {
				common.DeleteDeployment(data)
			} else {
				fmt.Printf("Deployment %q not exists\n", data.Metadata.Name)
			}
			break
		case "Service":
			if common.ServiceExistJudge(data) {
				common.DeleteService(data)
			} else {
				fmt.Printf("Service %q not exists\n", data.Metadata.Name)
			}
			break
		}
		break
	case "list":
		switch data.Kind {
		case "deployment":
			strback += common.DeploymentList(data)
			break
		case "service":
			strback += common.ServiceList(data)
			break

		}
		break

	}
	server.socket.Send(strback, 0)
}
