package storagenode

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"context"
)

type Server struct {
	store ShardStore
}

type ShardStore interface{
	Put(ctx context.Context, data io.Reader, size int64)(string,error)
	Get(ctx context.Context, cid string) ([]byte,error)
	// Pin(ctx context.Context, cid string) error
}

func NewServer() *Server {
	return &Server{store: &FileStore{}}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.readLoop(conn)
	}

}

func (s *Server) readLoop(conn net.Conn) error {
	buf := new(bytes.Buffer)
	fmt.Println("Reading data")
	for {
		var size int64
		var shardInd int32
		binary.Read(conn,binary.LittleEndian,&shardInd)
		binary.Read(conn, binary.LittleEndian, &size) // binary.Read() reads bytes from the connection and deserializes them into a value.
		// Both must use the same byte order (BigEndian/LittleEndian).
		
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		cid,err := s.store.Put(context.Background(),buf,size)
		if err!=nil{
			fmt.Println("Error is storing in shard store:",err)
		}
		fmt.Print(buf.Bytes())
		fmt.Println("Number of bytes:", n)
		fmt.Println("CID from IPFS:",cid)
	}

}

