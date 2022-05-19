package main

import (
	"flag"
	"k8s/common"
)

var realip string

func main() {
	flag.StringVar(&realip, "rip", "127.0.0.1", "Input the server ip,default 127.0.0.1")
	flag.Parse()
	done := make(chan bool)
	server := common.NewServer(realip, 20000, 20001)
	go server.Start(done)
	<-done
}
