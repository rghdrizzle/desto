package metadata

import (
	//"context"
	"encoding/json"
	"os"
	"path/filepath"
	types "rghdrizzle/desto/types"
	"fmt"
)
type LocalStore struct{

}
func (l *LocalStore) StoreFileMetaData( file *types.FileMetadata) (string, error){
	storeDir := getStoreDir()
	fileName := file.FileID
	path := filepath.Join(storeDir,fileName)
	b , err := json.Marshal(file)
	if err!=nil{
		return "",err
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return "", err
	}

	return path,nil
}

func (l *LocalStore) GetFileMetaData( fileid string) (*types.FileMetadata, error){

	fileMetadata := &types.FileMetadata{}
	storeDir := getStoreDir()
	path := filepath.Join(storeDir, fileid)
	jsonBytes, err:= os.ReadFile(path)
	if err!=nil{
		return nil,err
	}
	err = json.Unmarshal(jsonBytes,&fileMetadata)
	if err!=nil{
		return nil,err
	}
	
	return fileMetadata,nil
}

func getStoreDir() string{
	home, err :=os.UserHomeDir()
	if err!=nil{
		fmt.Println("error finding home dir:",err)
	}
	storeDir := filepath.Join(home, "desto-files")
	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return err.Error()
	}
	return storeDir
}