package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"rghdrizzle/desto/encryption"
)
var KEK string = "asdasdasdasdasdasdasw123dfhuilmn" // placeholder, will create generate KEK once at initial setup for the client and store it in the device key system


func SendFile(size int) error {
	
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		log.Fatal(err)
	}
	enveloped := encryption.NewEnveloped(&encryption.LocalKMS{Key: []byte(KEK)})
	conn, err := net.Dial("tcp", ":3001")
	if err != nil {
		log.Fatal(err)
	}
	encryptedData , err := enveloped.Encrypt(context.Background(),file)
	if err!=nil{
		fmt.Println("Error in encrypting file")
	}
	fmt.Println("Encrypted DATA:",encryptedData)
	binary.Write(conn, binary.LittleEndian, int64(size)) // binary.Write() serializes a value into bytes and writes them to the connection which then can be used by the binary reader to fetch and store it in a variable. This is useful to have additional data types for the connection
	n, err := io.CopyN(conn, bytes.NewReader([]byte(encryptedData)), int64(size))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Written bytes:", n)
	decryptedData , err := enveloped.Decrypt(context.Background(),encryptedData)
	fmt.Println("decrypted data:",decryptedData)
	return nil
}
