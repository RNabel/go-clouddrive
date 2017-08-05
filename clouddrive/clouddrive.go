package clouddrive

import (
	"github.com/boltdb/bolt"
	"google.golang.org/api/drive/v3"
	"bytes"
	"encoding/gob"
	"fmt"
	"CloudDrive/types"
)

type CloudFile struct {
	db *bolt.DB
	drv drive.Service
}

func NewCloudFileFromByte(in []byte) *CloudFile {
	cf := CloudFile{}

	buf := bytes.Buffer{}
	buf.Write(in)

	decoder := gob.NewDecoder(&buf)
	err := decoder.Decode(&cf)
	if err != nil {
		fmt.Println("failed to gob Decode", err)
	}

	return &cf
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