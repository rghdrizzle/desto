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
		client.ReceiveFile("d6b7607f-bd31-4ac0-8ea4-75abe33acd0f")
	}()
	fmt.Println("server started")
	server := sn.NewServer()
	server.Start()
}