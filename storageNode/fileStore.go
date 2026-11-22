package storagenode

import (
	"context"
	"fmt"
	"io"
	"os"
	"crypto/sha1"
	"encoding/hex"
	"bytes"
	"path/filepath"
)

type FileStore struct {
}

func (f *FileStore) Put(ctx context.Context,data io.Reader,size int64) (string, error){
	buf := new(bytes.Buffer)
	_, err := io.CopyN(buf, data, int64(size))
	if err!=nil{
		fmt.Println("error in copying data to flle:",err)
	}
	// filename := GetHashFromFileForFilename(buf.String()[:12])
	// file,err:= os.Create(filename)
	// if err!=nil{
	// 	fmt.Println("error in creating flle:",err)
	// }
	// n, err := io.CopyN(file, buf, int64(size))
	// fmt.Println("File created with data with written bytes:",n)
	// return filename,nil

	filename := GetHashFromFileForFilename(buf.String()[:12])
	storeDir := getStoreDir()
	path := filepath.Join(storeDir, filename)
	file , err := os.OpenFile(path,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,0644)
	if err!=nil{
		return "",err
	}
	defer file.Close()
	_,err= file.Write(buf.Bytes())
	if err!=nil{
		return "",err
	}
	if err := file.Sync();err!=nil{
		return "",err
	}
	// if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
	// 	return "", err
	// }
	return filename, nil
}


func GetHashFromFileForFilename(key string) string {
	hash := sha1.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (f *FileStore) Get(ctx context.Context, cid string) ([]byte, error) {

	filename := cid
	storeDir := getStoreDir()
	path := filepath.Join(storeDir, filename)
	return os.ReadFile(path)
}

func getStoreDir() string{
	home, err :=os.UserHomeDir()
	if err!=nil{
		fmt.Println("error finding home dir:",err)
	}
	storeDir := filepath.Join(home, "desto-shards")
	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return err.Error()
	}
	return storeDir
}