package main

import (
	"flag"
	"fmt"
	"github.com/pebbe/zmq4"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

func startClient(port int, done chan bool, msg []byte) {
	// REQ  表示client端
	socket, _ := zmq4.NewSocket(zmq4.REQ)
	//绑定端口，指定传输层协议
	socket.Connect("tcp://127.0.0.1:" + strconv.Itoa(port))
	fmt.Printf("connect to server\n")
	defer socket.Close()

	socket.SendBytes(msg, zmq4.DONTWAIT)
	socket.RecvBytes(0)
	//fmt.Printf("%s", rec)
	done <- true
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Your input args is too less,it should be at least 4 paramemters just like\n 'kubectl create -f xxx.yaml'")
		return
	}
	var file string
	flag.StringVar(&file, "f", "examples/tomcat-dp.yaml", "Input your yaml file")
	flag.Parse()

	var init func(op string)
	init = func(op string) {
		u, _ := user.Current()
		user := []byte("user: " + u.Name)
		operation := []byte("\noperation: " + op)
		done := make(chan bool)
		buff, _ := ioutil.ReadFile(file)
		buff = append(buff, user...)
		buff = append(buff, operation...)
		go startClient(20000, done, buff)
		<-done
	}

	switch os.Args[1] {
	case "apply":
		switch os.Args[2] {
		case "-f":
			_, err := os.Stat(os.Args[3])
			if err != nil {
				fmt.Println("Your input file is not exists")
				return
			}
			init("apply")
		case "-i":

		default:
			fmt.Println("input error -> go for help")
		}
	case "delete":
		switch os.Args[2] {
		case "-f":
			_, err := os.Stat(os.Args[3])
			if err != nil {
				fmt.Println("Your input file is not exists")
				return
			}
			fmt.Println(file)
			init("delete")
		case "-i":

		default:
			fmt.Println("input error -> go for help")
		}
	default:
		fmt.Println("operation err go for help")
	}

}
