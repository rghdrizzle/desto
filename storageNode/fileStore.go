package storagenode

import (
	"context"
	"fmt"
	"io"
	"os"
	"crypto/sha1"
	"encoding/hex"
	"bytes"
)

type FileStore struct {
}

func (f *FileStore) Put(ctx context.Context,data io.Reader,size int64) (string, error){
	buf := new(bytes.Buffer)
	_, err := io.CopyN(buf, data, int64(size))
	if err!=nil{
		fmt.Println("error in copying data to flle:",err)
	}
	filename := GetHashFromFileForFilename(buf.String()[:12])
	file,err:= os.Create(filename)
	if err!=nil{
		fmt.Println("error in creating flle:",err)
	}
	n, err := io.CopyN(file, buf, int64(size))
	fmt.Println("File created with data with written bytes:",n)
	return filename,nil
}


func GetHashFromFileForFilename( key string ) string{
	hash := sha1.Sum([]byte(key)) // [20]bytes => to a slice([]byte) by doing => [:]
	hashString := hex.EncodeToString(hash[:])
	fmt.Println("Hash:",hashString)
	return hashString
}

func (f *FileStore) Get(ctx context.Context,cid string) ([]byte,error){
	file , err := os.Open(cid)
	if err!=nil{
		fmt.Println("Error reading file:",err)
	}
	buf := new(bytes.Buffer)
	n,err:= io.Copy(buf,file)
	if err!=nil{
		fmt.Println("Error copying file data to buffer:",err)
	}
	fmt.Println("buffer created with data with written bytes:",n)
	return buf.Bytes(),nil
}