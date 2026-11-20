package client

import (
	"bytes"
	"context"
	"time"

	//"crypto/rand"
	//"encoding/binary"
	"fmt"
	//"io"
	"log"
	"os"

	//"net"
	metadata "rghdrizzle/desto/client/metadata"
	"rghdrizzle/desto/encryption"
	transport "rghdrizzle/desto/transport"
	types "rghdrizzle/desto/types"
	"github.com/google/uuid"
	"github.com/klauspost/reedsolomon"
	"google.golang.org/grpc"
)

type FileMetaDataStore interface {
	StoreFileMetaData( file types.FileMetadata) (string, error)
	GetFileMetaData( fileid string) (types.FileMetadata, error)
}

var KEK string = "asdasdasdasdasdasdasw123dfhuilmn" // placeholder, will create generate KEK once at initial setup for the client and store it in the device key system

type ErasureEncoder struct {
	dataShards   int
	parityShards int
}

var store *metadata.LocalStore = &metadata.LocalStore{}

// TODO: Refactor code to make sure to handle fetch files from the server and manage verification of data after retrieval
func SendFile(size int) error {
	fileMetadata := &types.FileMetadata{Shards: []types.ShardMetadata{}}
	file := []byte("YEA THIS WORKS!!!asdadasd") //make([]byte, size)
	enveloped := encryption.NewEnveloped(&encryption.LocalKMS{Key: []byte(KEK)})
	// grpc
	conn1, err := grpc.NewClient("localhost:3002", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("client active listening at %d\n", 3002)
	defer conn1.Close()
	c := transport.NewDestoClient(conn1)
	encryptedData, err := enveloped.Encrypt(context.Background(), file)
	sizeOfEncryptedData := len(encryptedData) // TODO: should send to the server when sending the file since we need it for decoding the shards and decryption
	if err != nil {
		fmt.Println("Error in encrypting file")
	}
	fmt.Println("Encrypted DATA:", encryptedData)
	// erasure coding of this data
	erasureEncoder := &ErasureEncoder{dataShards: 4, parityShards: 2}
	shards, err := erasureEncoder.erasureEncoding(encryptedData)
	for i, shard := range shards {
		cid, err := c.SendFile(context.Background(), &transport.File{Data: shard, Size: int64(len(shard)), Index: int64(i)})
		checkErr(err)
		fileMetadata.Shards = append(fileMetadata.Shards,types.ShardMetadata{CID: cid.Cid,Index: i})
		fmt.Println("cid returned:", cid.Cid)
	}
	fileMetadata.CreatedAt = time.Now()
	fileMetadata.FileID = uuid.NewString()
	fileMetadata.Name = "Test"
	fileMetadata.Size = int64(sizeOfEncryptedData)
	// create metadata file
	path , err := store.StoreFileMetaData(fileMetadata)
	checkErr(err)
	fmt.Println("Metadata path:",path)

	// // testing with data from here since I did not implement recieving data from the server yet, so we will use the shards data from this function itself for verification of business logic
	// placeholderDATAShardsFromServer := shards
	// dataFromShard, err := erasureEncoder.erasureDecoding(placeholderDATAShardsFromServer, sizeOfEncryptedData) // replace size with size from the server
	// checkErr(err)
	// // TODO: fix the bug related to mismatch in size since we do not know the actual file size when decoding which is important when it comes to joining the shard
	// decryptedData, err := enveloped.Decrypt(context.Background(), dataFromShard)
	// fmt.Printf("Decrypted bytes: %v\n", string(decryptedData))
	return nil
}

func ReceiveFile(fileid string) error {
	// TODO: fetch details from the manifest like file size, number of shards, shard index and so on

	enveloped := encryption.NewEnveloped(&encryption.LocalKMS{Key: []byte(KEK)})
	erasureEncoder := &ErasureEncoder{dataShards: 4, parityShards: 2}

	// grpc
	conn1, err := grpc.NewClient("localhost:3002", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("client active listening at %d\n", 3002)
	defer conn1.Close()
	c := transport.NewDestoClient(conn1)
	fileMetadata ,err:= store.GetFileMetaData(fileid)
	checkErr(err)
	shards := [][]byte{}
	for _ ,shard:= range fileMetadata.Shards{
		shard, err := c.RetrieveFile(context.Background(), &transport.Cid{Cid: shard.CID})
		checkErr(err)
		shards = append(shards,shard.Data)
	}

	fmt.Println("shards length retrieved", len(shards))
	dataFromShard, err := erasureEncoder.erasureDecoding(shards, int(fileMetadata.Size)) // replace size with size from the server
	checkErr(err)
	// TODO: fix the bug related to mismatch in size since we do not know the actual file size when decoding which is important when it comes to joining the shard
	decryptedData, err := enveloped.Decrypt(context.Background(), dataFromShard)
	fmt.Printf("Decrypted bytes: %v\n", string(decryptedData))
	return nil
}

func (e *ErasureEncoder) erasureEncoding(data string) ([][]byte, error) {
	if (e.dataShards + e.parityShards) > 256 {
		log.Fatalf("Error: sum of data and parity shards cannot exceed 256\n")
	}
	enc, err := reedsolomon.New(e.dataShards, e.parityShards)
	checkErr(err)

	shards, err := enc.Split([]byte(data))
	checkErr(err)

	err = enc.Encode(shards)
	checkErr(err)

	// shards with parity
	return shards, nil

}

func (e *ErasureEncoder) erasureDecoding(shards [][]byte, filesize int) (string, error) {
	if (e.dataShards + e.parityShards) > 256 {
		log.Fatalf("Error: sum of data and parity shards cannot exceed 256\n")
	}
	enc, err := reedsolomon.New(e.dataShards, e.parityShards)
	checkErr(err)
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
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
	err = enc.Join(buf, shards, filesize)
	checkErr(err)

	data := buf.Bytes()
	fmt.Printf("\n Encrypted data from server: %s\n", string(data))
	return string(data), err

}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
