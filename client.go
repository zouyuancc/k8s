package main

import (
	"flag"
	"fmt"
	"github.com/pebbe/zmq4"
	"io/ioutil"
	"os"
	user2 "os/user"
	"strconv"
)

var file string
var operation string
var serverIp string
var serverPort int
var doneClient chan bool

func init() {
	flag.StringVar(&file, "f", "examples/tomcat-dp.yaml", "Input your yaml file")
	flag.StringVar(&operation, "op", "apply", "Input you operation,like \"apply,delete and so on\"")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Input you remote ip address default 127.0.0.1")
	flag.IntVar(&serverPort, "port", 20000, "Input the remote server port to connect,default value 20000")
}

type Client struct {
	Ip   string
	Port int

	socket *zmq4.Socket
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		Ip:   ip,
		Port: port,
	}
	socket, _ := zmq4.NewSocket(zmq4.REQ)
	client.socket = socket
	err := client.socket.Connect("tcp://" + ip + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("client connect err:", err)
		return nil
	}
	return client
}

func (client *Client) FileMethod() {
	defer client.socket.Close()
	u, _ := user2.Current()
	user := []byte("user: " + u.Name)
	operation := []byte("\noperation: " + operation)
	buff, _ := ioutil.ReadFile(file)
	buff = append(buff, user...)
	buff = append(buff, operation...)

	client.socket.SendBytes(buff, zmq4.DONTWAIT)
	data, _ := client.socket.RecvBytes(0)
	fmt.Println(string(data))
	doneClient <- true
}

func (client *Client) CliMethod() {
	defer client.socket.Close()

}

func (client *Client) JudgeOption() {
	switch operation {
	case "apply":
		switch file {
		case "":

		default:
			_, err := os.Stat(file)
			if err != nil {
				fmt.Println("Your input file is not exists")
				doneClient <- false
				return
			}
			client.FileMethod()
		}
	case "delete":
		switch file {
		case "":

		default:
			_, err := os.Stat(file)
			if err != nil {
				fmt.Println("Your input file is not exists")
				doneClient <- false
				return
			}
			fmt.Println(file)
			client.FileMethod()
		}
	default:
		fmt.Println("operation err go for help")
	}
}

func (client *Client) Response() {
	resp, _ := client.socket.Recv(0)
	fmt.Println(resp)
	client.socket.Send("recvived", 0)
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	doneClient = make(chan bool, 1)
	client.JudgeOption()
	go client.Response()
	<-doneClient
}
