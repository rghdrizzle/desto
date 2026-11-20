package types
import (
	"time"
)

type ShardMetadata struct{
	Index int   `json:"index"`
	CID    string   `json:"cid"`
    //Nodes  []string `json:"nodes"`
    Size   int64    `json:"size"`
}
type FileMetadata struct{
	Version   int       `json:"version"`
    FileID    string    `json:"file_id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    Shards []ShardMetadata `json:"shards"`
}