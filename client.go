package main

import (
	"flag"
	"fmt"
	"github.com/pebbe/zmq4"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"k8s/cores"
	"net"
	"os"
	user2 "os/user"
	"strconv"
)

var file string
var operation string
var serverIp string
var serverPort1 int
var serverPort2 int
var doneClient chan bool
var kind string
var namespace string

func init() {
	flag.StringVar(&file, "f", "examples/tomcat-dp.yaml", "Input your yaml file")
	flag.StringVar(&operation, "op", "list", "Input you operation,like \"apply,delete and so on\"")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Input you remote ip address default 127.0.0.1")
	flag.IntVar(&serverPort1, "port1", 20000, "Input the remote server port to connect,default value 20000")
	flag.IntVar(&serverPort2, "port2", 20001, "Input the remote server port to connect,default value 20000")
	flag.StringVar(&kind, "kind", "service", "Input the kind of resources,default value 20000")
	flag.StringVar(&namespace, "nm", "default", "Input the namespace of resources,default value default")
}

//task client
type Client struct {
	Ip   string
	Port int

	data cores.Yaml

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
	data, _ := client.socket.Recv(0)
	fmt.Println(data)
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
	case "list":
		client.data.Kind = kind
		client.data.Metadata.Namespace = namespace
		u, _ := user2.Current()
		client.data.Operation = operation
		client.data.User = u.Name
		buf, _ := yaml.Marshal(client.data)
		client.socket.SendBytes(buf, zmq4.DONTWAIT)
		data, _ := client.socket.RecvBytes(0)
		fmt.Println(string(data))
		doneClient <- true

	default:
		fmt.Println("operation err go for help")
	}
}

func (client *Client) Response() {
	resp, _ := client.socket.Recv(0)
	fmt.Println(resp)
	client.socket.Send("recvived", 0)
}

func InitTask() {
	client := NewClient(serverIp, serverPort1)
	doneClient = make(chan bool, 1)
	client.JudgeOption()
	<-doneClient
}

//user client
type Cclient struct {
	Ip   string
	Port int

	conn   net.Conn
	symbol int
}

func NewCclient(ip string, port int) *Cclient {
	client := &Cclient{
		Ip:   ip,
		Port: port,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("client conn err:", err)
		return nil
	}
	client.conn = conn
	go client.Response()
	return client
}

func (client Cclient) menu() {
	fmt.Println("----------帮助----------")
	fmt.Println("0.查看帮助文档")
	fmt.Println("1.使用当前传入配置（参数+文件）配置集群任务")
	fmt.Println("2.更改任务启动使用的yaml配置文件")
	fmt.Println("3.查询在线用户")
	fmt.Println("4.查看指定用户任务记录")
	fmt.Println("11.退出")
	fmt.Println("----------------------")
	fmt.Println()
}

func (client *Cclient) Response() {
	io.Copy(os.Stdout, client.conn)
	for {
		buf := make([]byte, 4096)
		n, _ := client.conn.Read(buf)
		if n == 0 {
			fmt.Println("服务端断开连接")
			break
		}
		fmt.Println(buf)
	}

}

func (client Cclient) view_all_users() {
	msg := "who\n"
	client.conn.Write([]byte(msg))
}

func (client Cclient) Run() {
	client.menu()
	var sign int
	fmt.Printf("请输入：")
	fmt.Scan(&sign)

	for sign != 11 {
		client.symbol = sign
		switch client.symbol {
		case 0:
			client.menu()
			break
		case 1:
			fmt.Println("选择使用使用当前配置（参数+文件）配置集群任务")
			InitTask()
			break
		case 2:
			fmt.Println("更改任务启动使用的yaml配置文件")
			break
		case 3:
			fmt.Println("查询在线用户")
			client.view_all_users()
			break
		case 4:
			fmt.Println("查看指定用户任务记录")
			break
		case 11:
			return
		}

		sign = 0
		fmt.Printf("请输入：")
		fmt.Scan(&sign)
	}
}

func main() {
	flag.Parse()
	client := NewCclient(serverIp, serverPort2)
	client.Run()
}
