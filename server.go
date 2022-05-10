package main

import (
	"encoding/json"
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	yaml "gopkg.in/yaml.v2"
	"k8s/common"
	"k8s/cores"
	"log"
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
		go server.Parseargs([]byte(resp))
		fmt.Println("ok,receive")
		server.socket.Send("We have received your  task request!", 0)
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
	data := new(cores.Yaml)
	yaml.Unmarshal(resp, data)
	lins, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data:\t", string(lins))

	switch data.Operation {
	case "apply":
		switch data.Kind {
		case "Deployment":
			if common.DeploymentExistJudge(data) {
				common.UpdateDeployment(data)
			} else {
				common.CreateDeployment(data)
			}
		case "Service":
			if common.ServiceExistJudge(data) {
				common.UpdateService(data)
			} else {
				common.CreateService(data)
			}
		}
	case "delete":
		switch data.Kind {
		case "Deployment":
			if common.DeploymentExistJudge(data) {
				common.DeleteDeployment(data)
			} else {
				fmt.Printf("Deployment %q not exists\n", data.Metadata.Name)
			}
		case "Service":
			if common.ServiceExistJudge(data) {
				common.DeleteService(data)
			} else {
				fmt.Printf("Service %q not exists\n", data.Metadata.Name)
			}
		}
	case "list":
		switch data.Kind {
		case "Deployment":

		case "Service":

		}

	}

}
