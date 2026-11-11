package main

import (
	"fmt"
	"os"
	client "rghdrizzle/desto/client"
	sn "rghdrizzle/desto/storageNode"
	"time"
)

func main(){
	go func(){
		time.Sleep(4*time.Second)
		client.SendFile(200)
		os.Exit(0)
	}()
	fmt.Println("server started")
	server := sn.NewServer()
	server.Start()
}