package storagenode

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
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
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size) // binary.Read() reads bytes from the connection and deserializes them into a value.
		// Both must use the same byte order (BigEndian/LittleEndian).
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(buf.Bytes())
		fmt.Println("Number of bytes:", n)
	}

}
