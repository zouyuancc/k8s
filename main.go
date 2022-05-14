package main

import "k8s/common"

func main() {
	done := make(chan bool)
	server := common.NewServer("192.168.163.111", 20000, 20001)
	go server.Start(done)
	<-done
}
