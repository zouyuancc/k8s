package main

import (
	zmq "github.com/pebbe/zmq4"
	yaml "gopkg.in/yaml.v2"
	"k8s/pkg/parse"
	yaml_define "k8s/pkg/stru"
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
	conf := new(yaml_define.Yaml)
	yaml.Unmarshal(resp, conf)
	//fmt.Println(conf)
	parse.OperateSource(conf)
}
