package main

import (
	//"fmt"
	//"os"
	//client "rghdrizzle/desto/client"
	sn "rghdrizzle/desto/storageNode"
	//"time"
)

func main(){
	ready := make(chan struct{})
	server := sn.NewServer()
	go server.Start(ready)
	<-ready
}