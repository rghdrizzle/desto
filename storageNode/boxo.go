package storagenode

import (
	//"github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
	"fmt"
)


func Test(shardData []byte)error {
	prefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}

	block, err := blocks.NewBlockWithPrefix(shardData, prefix)
	if err != nil {
		return err
	}
	fmt.Println(block.Cid())
	return nil
}
