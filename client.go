package main

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"io/ioutil"
	"strconv"
)

func startClient(port int, done chan bool, msg []byte) {
	// REQ  表示client端
	socket, _ := zmq4.NewSocket(zmq4.REQ)
	//绑定端口，指定传输层协议
	socket.Connect("tcp://127.0.0.1:" + strconv.Itoa(port))
	fmt.Printf("connect to server")
	defer socket.Close()

	socket.SendBytes(msg, zmq4.DONTWAIT)
	socket.RecvBytes(0)
	//fmt.Printf("%s", rec)
	done <- true
}

func main() {
	done := make(chan bool)
	buff, _ := ioutil.ReadFile("/home/zy/GolandProjects/k8s/examples/test.yaml")
	go startClient(20000, done, buff)
	<-done
}
