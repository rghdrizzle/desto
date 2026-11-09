package main

import (
	"fmt"
	sn "rghdrizzle/desto/storageNode"
	client "rghdrizzle/desto/client"
	"time"
)

func main(){
	go func(){
		time.Sleep(4*time.Second)
		client.SendFile(2000)
	}()
	fmt.Println("server started")
	server := sn.NewServer()
	server.Start()
}