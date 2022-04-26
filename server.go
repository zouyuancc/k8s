package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"strconv"
	parsedp "test/pkg"
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
	conf := new(parsedp.Yaml)
	yaml.Unmarshal(resp, conf)
	//fmt.Println(conf)
	//parsedp.CreateSource(conf)
	data, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data:\t", string(data))
}
