package tests

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	//"os"
	"rghdrizzle/desto/client"
	sn "rghdrizzle/desto/storageNode"
	"testing"
)

func TestSendandReceiveFile(t *testing.T) {

	ready := make(chan struct{})
	server := sn.NewServer()
	go server.Start(ready)
	<-ready

	image1,_:=os.ReadFile("pixil-frame-0(5).png")
	testFiles := [][]byte{
		[]byte("Testing string 1"),
		[]byte("Testing string 2"),
		image1,
	}
	for i := range testFiles {
		file := testFiles[i]
		fmt.Println("Testing starts: ")
		fileid, err := client.SendFile(file,("test"+strconv.Itoa((i))))
		if err != nil {
			t.Fatal("Error sending file to the server:", err)
		}
		recievedData, err := client.ReceiveFile(fileid)
		if err != nil {
			t.Fatal("Error fetching file from the server:", err)
		}
		if !bytes.Equal(recievedData, file) {
			t.Fatal("File fetched is corrupted or not the same file uploaded")
		}
		t.Log("Files are intact and coherent after upload and retrieval")
	}
}
