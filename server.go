package main

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	yaml "gopkg.in/yaml.v2"
	"k8s/common"
	"k8s/cores"
	"log"
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
		
	}

}
