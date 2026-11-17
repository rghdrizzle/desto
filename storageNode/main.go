package storagenode

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"rghdrizzle/desto/transport"

	"google.golang.org/grpc"
)

type Server struct {
	store ShardStore
	transport.UnimplementedDestoServer
}

type ShardStore interface {
	Put(ctx context.Context, data io.Reader, size int64) (string, error)
	Get(ctx context.Context, cid string) ([]byte, error)
	// Pin(ctx context.Context, cid string) error
}

func NewServer() *Server {
	return &Server{store: &FileStore{}}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", ":3002")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	transport.RegisterDestoServer(server, s)
	fmt.Println("Registered grpc server")
	if err := server.Serve(ln); err != nil {
		log.Fatal(err)
	}

}

func (s *Server) SendFile(context context.Context, file *transport.File) (*transport.Cid, error) {
	fmt.Println("Sending file")
	buf := new(bytes.Buffer)
	// var size int64
	// var shardInd int32
	// binary.Read(conn, binary.LittleEndian, &shardInd)
	// binary.Read(conn, binary.LittleEndian, &size) // binary.Read() reads bytes from the connection and deserializes them into a value.
	// // Both must use the same byte order (BigEndian/LittleEndian).
	reader := bytes.NewReader(file.Data)
	n, err := io.CopyN(buf, reader, file.Size)
	if err != nil {
		log.Fatal(err)
	}
	cid, err := s.store.Put(context, buf, file.Size)
	if err != nil {
		fmt.Println("Error is storing in shard store:", err)
	}
	fmt.Print(buf.Bytes())
	fmt.Println("Number of bytes:", n)
	fmt.Println("CID from IPFS:", cid)
	return &transport.Cid{Cid: cid},nil
}

func (s *Server) readLoop(conn net.Conn) error {
	buf := new(bytes.Buffer)
	fmt.Println("Reading data")
	for {
		var size int64
		var shardInd int32
		binary.Read(conn, binary.LittleEndian, &shardInd)
		binary.Read(conn, binary.LittleEndian, &size) // binary.Read() reads bytes from the connection and deserializes them into a value.
		// Both must use the same byte order (BigEndian/LittleEndian).

		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		cid, err := s.store.Put(context.Background(), buf, size)
		if err != nil {
			fmt.Println("Error is storing in shard store:", err)
		}
		fmt.Print(buf.Bytes())
		fmt.Println("Number of bytes:", n)
		fmt.Println("CID from IPFS:", cid)
	}

}
