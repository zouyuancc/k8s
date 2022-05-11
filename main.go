package main

import "k8s/common"

func main() {
	done := make(chan bool)
	server := common.NewServer("127.0.0.1", 20000, 20001)
	go server.Start(done)
	<-done
}
