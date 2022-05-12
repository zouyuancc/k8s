package common

import (
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

var ip string
var port1 int
var port2 int

func init() {
	flag.StringVar(&ip, "ip", "192.168.163.111", "Input the server ip,default 127.0.0.1")
	flag.IntVar(&port1, "port1", 20000, "Input the server port,default 20000")
	flag.IntVar(&port2, "port2", 20001, "Input the server port,default 20000")
}

type Server struct {
	Ip    string
	Port1 int
	Port2 int

	//任务级用户表
	OnlineMap map[string]*History
	Maplock   sync.RWMutex

	//任务级用户表
	COnlineMap map[string]*CUser
	CMaplock   sync.RWMutex
}

func NewServer(ip string, port1, port2 int) *Server {
	server := &Server{
		Ip:         ip,
		Port1:      port1,
		Port2:      port2,
		OnlineMap:  make(map[string]*History),
		COnlineMap: make(map[string]*CUser),
	}
	return server
}

func (server *Server) Handle_chat(conn net.Conn) {
	//创建用户级别：
	cuser := NewCUser(conn, server)

	cuser.Online()

	//监听用户是否活跃的管道
	isLive := make(chan bool)

	//接收用户发来的消息，并处理
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				cuser.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			//提取用户消息
			msg := string(buf[:n-1])
			cuser.Domessage(msg)
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-isLive:
			fmt.Println(cuser.Name, "refresh")
		case <-time.After(time.Hour * 24):
			cuser.SendMsg([]byte("超时退出\n"))
			close(cuser.C)
			conn.Close()
			return
		}
	}

}

func (server *Server) Chat_with_client() {
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port2))
	if err != nil {
		fmt.Println("net listen err on port2:", err)
		return
	}
	defer Listener.Close()
	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("listener accept err on port2", err)
			return
		}

		go server.Handle_chat(conn)
	}
}

func (server *Server) Handle_task(socket zmq.Socket) {

}

func (server *Server) Start(done chan bool) {
	go server.Chat_with_client()
	socket, _ := zmq.NewSocket(zmq.REP)
	user := NewUser(socket, server)
	for {
		user.Socket.Bind("tcp://" + ip + ":" + strconv.Itoa(port1))
		//Recv 和 Send必须交替进行
		resp, _ := user.Socket.Recv(0)
		if resp != "" {
			go user.Parseargs([]byte(resp))
		}
	}
	done <- true

}
