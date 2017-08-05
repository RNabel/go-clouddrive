package types

type FileState struct {
	DownloadTimestamp int64
	ExpiryDate int64
	Filepath string
	// TODO More file information.
}

type File interface {
	ToBytes() []byte
}
