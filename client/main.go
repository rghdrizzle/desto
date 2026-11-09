package client

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func SendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.Dial("tcp", ":3001")
	if err != nil {
		log.Fatal(err)
	}
	binary.Write(conn, binary.LittleEndian, int64(size)) // binary.Write() serializes a value into bytes and writes them to the connection which then can be used by the binary reader to fetch and store it in a variable. This is useful to have additional data types for the connection
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Written bytes:", n)
	return nil
}
