package types

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"google.golang.org/api/drive/v3"
)

type FileState struct {
	DownloadTimestamp int64
	ExpiryDate int64
	Filepath string
	// TODO More file information.
}

// Export types.
type CloudFile struct {
	GoogleId string
	Name string
	Parents []string
}

func NewCloudFile(file *drive.File) CloudFile {
	return CloudFile{file.Id, file.Name, file.Parents}
}

func NewCloudFileFromByte(in []byte) CloudFile {
	cf := CloudFile{}

	buf := bytes.Buffer{}
	buf.Write(in)

	decoder := gob.NewDecoder(&buf)
	err := decoder.Decode(&cf)
	if err != nil {
		fmt.Println("failed to gob Decode", err)
	}

	return cf
}

func (cf *CloudFile) ToBytes() []byte {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(cf)
	if err != nil {
		fmt.Println("failed to gob Encode", err)
	}
	return b.Bytes()
}