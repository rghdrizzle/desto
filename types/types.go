package types
import (
	"time"
)

type ShardMetadata struct{
	Index int   `json:"index"`
	CID    string   `json:"cid"`
    //Nodes  []string `json:"nodes"`
}
type FileMetadata struct{
	Version   int       `json:"version"`
    FileID    string    `json:"file_id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
	Size   int64    `json:"size"`
    Shards []ShardMetadata `json:"shards"`
}