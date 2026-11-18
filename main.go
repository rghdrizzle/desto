package main

import (
	"fmt"
	//"os"
	client "rghdrizzle/desto/client"
	sn "rghdrizzle/desto/storageNode"
	"time"
)

func main(){
	go func(){
		time.Sleep(4*time.Second)
		client.SendFile(200)
		//os.Exit(0)
		time.Sleep(4*time.Second)
		client.ReceiveFile("d229b30813f079bb42afc5e73e9d4d85cf0a3cb9")
	}()
	fmt.Println("server started")
	server := sn.NewServer()
	server.Start()
}