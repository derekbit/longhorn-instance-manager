package client

type DiskInfo struct {
	ID          string
	UUID        string
	Path        string
	Type        string
	TotalSize   int64
	FreeSize    int64
	TotalBlocks int64
	FreeBlocks  int64
	BlockSize   int64
	ClusterSize int64
	Readonly    bool
}

type ReplicaInfo struct {
	Name        string
	UUID        string
	BdevName    string
	LvstoreUUID string

	TotalSize   int64
	TotalBlocks int64
	BlockSize   int64

	ThinProvision bool

	State string
}