package client

import (
	"bytes"
	"context"
	//"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"log"
	"net"
	"rghdrizzle/desto/encryption"
	"github.com/klauspost/reedsolomon"
)
var KEK string = "asdasdasdasdasdasdasw123dfhuilmn" // placeholder, will create generate KEK once at initial setup for the client and store it in the device key system

type ErasureEncoder struct{
	dataShards int
	parityShards int
}
// TODO: Refactor code to make sure to handle fetch files from the server and manage verification of data after retrieval
func SendFile(size int) error {
	
	file := []byte("YEA THIS WORKS!!!asdadasd")//make([]byte, size)
	// _, err := io.ReadFull(rand.Reader, file)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	enveloped := encryption.NewEnveloped(&encryption.LocalKMS{Key: []byte(KEK)})
	conn, err := net.Dial("tcp", ":3001")
	if err != nil {
		log.Fatal(err)
	}
	encryptedData , err := enveloped.Encrypt(context.Background(),file)
	sizeOfEncryptedData := len(encryptedData) // TODO: should send to the server when sending the file since we need it for decoding the shards and decryption
	if err!=nil{
		fmt.Println("Error in encrypting file")
	}
	fmt.Println("Encrypted DATA:",encryptedData)
	// erasure coding of this data
	erasureEncoder:= &ErasureEncoder{dataShards: 4,parityShards: 2}
	shards ,err := erasureEncoder.erasureEncoding(encryptedData)
	for i, shard := range shards{
		// we need index for consistent reconstruction during decoding
		if err := binary.Write(conn, binary.LittleEndian, uint32(i)); err != nil {
            return err
        }
        // Send shard size.
        if err := binary.Write(conn, binary.LittleEndian, uint64(len(shard))); err != nil {
            return err
        }
        n, err := io.CopyN(conn, bytes.NewReader(shard),int64(len(shard)))
		if err != nil {
            return err
        }
		fmt.Println("Written bytes:", n)
	} 
	
	// testing with data from here since I did not implement recieving data from the server yet, so we will use the shards data from this function itself for verification of business logic
	placeholderDATAShardsFromServer := shards
	dataFromShard,err:=erasureEncoder.erasureDecoding(placeholderDATAShardsFromServer,sizeOfEncryptedData) // replace size with size from the server
	checkErr(err)
	// TODO: fix the bug related to mismatch in size since we do not know the actual file size when decoding which is important when it comes to joining the shard
	decryptedData , err := enveloped.Decrypt(context.Background(),dataFromShard)
	fmt.Printf("Decrypted bytes: %v\n", string(decryptedData))
	return nil
}

func receiveFile() error{
	
	
	return nil
}

func (e * ErasureEncoder) erasureEncoding(data string) ([][]byte, error){
	if (e.dataShards + e.parityShards) > 256 {
		log.Fatalf("Error: sum of data and parity shards cannot exceed 256\n")
	}
	enc, err := reedsolomon.New(e.dataShards, e.parityShards)
	checkErr(err)

	shards,err := enc.Split([]byte(data))
	checkErr(err)

	err = enc.Encode(shards)
	checkErr(err)

	// shards with parity
	return shards,nil



}


func (e * ErasureEncoder) erasureDecoding(shards [][]byte, filesize int) (string, error){
	if (e.dataShards + e.parityShards) > 256 {
		log.Fatalf("Error: sum of data and parity shards cannot exceed 256\n")
	}
	enc, err := reedsolomon.New(e.dataShards, e.parityShards)
	checkErr(err)
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	}else {
		fmt.Println("Verification failed. Reconstructing data")
		err = enc.Reconstruct(shards)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			os.Exit(1)
		}
		checkErr(err)
	}
	buf := new(bytes.Buffer)
	err = enc.Join(buf,shards,filesize)
	checkErr(err)

	data := buf.Bytes()
	fmt.Printf("\n Encrypted data from server: %s\n",string(data))
	return string(data),err



}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}