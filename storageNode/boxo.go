package storagenode

import (
	//"github.com/ipfs/boxo/blockstore"
	"bytes"
	"fmt"
	"log"

	"context"
	"io"

	"github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ds-pebble"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
)

var prefix = cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}

type boxo struct{
	ds datastore.Batching
	bs blockstore.Blockstore
}

func NewBoxoStore() *boxo{
	ds ,err := pebbleds.NewDatastore(GetStoreDir("boxo"))
	bs := blockstore.NewBlockstore(ds)
	if err != nil {
		log.Fatal(err)
	}
	return &boxo{ds: ds,bs: bs}
}

func (b *boxo) Put(ctx context.Context, data io.Reader, size int64) (string, error){
	buf := new(bytes.Buffer)
	_, err := io.CopyN(buf, data, size)
	if err!=nil{
		return "",err
	}
	block, err := blocks.NewBlockWithPrefix(buf.Bytes(), prefix)
	if err != nil {
		return "",err
	}
	fmt.Println(block.Cid())
	err = b.bs.Put(ctx,block)
	if err != nil {
		return "",err
	}
	return block.Cid().String(),nil
}

func (b *boxo) Get(ctx context.Context, id string) ([]byte, error){
	cid,err:= cid.Decode(id)
	if err!=nil{
		return nil,err
	}
	block ,err:= b.bs.Get(ctx,cid)
	if err!=nil{
		return nil,err
	}
	return block.RawData(), nil
}

func (b *boxo) Close() error {
    return b.ds.Close()
}