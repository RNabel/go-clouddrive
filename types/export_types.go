package types

import "google.golang.org/api/drive/v3"

type FileState struct {
	DownloadTimestamp int64
	ExpiryDate int64
	Filepath string
	// TODO More file information.
}

type File interface {
	ToBytes() []byte
	CopyGoogleFile(file *drive.File)
	FromBytes(in []byte)
	Parents() []string
	Name() string
	Id() string
}
