package files

import (
	"google.golang.org/api/drive/v3"
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
	"github.com/boltdb/bolt"

	"CloudDrive/database"
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
}

func NewCloudFile(file *drive.File) CloudFile {
	return CloudFile{file.Id}
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

func AddFilesToDB(fileChan chan CloudFile, db *bolt.DB, wg sync.WaitGroup) {
	defer wg.Done()

	// Create map from IDs to paths.


	var count = 0

	for cloudFile := range fileChan {
		count += 1
		serialised := cloudFile.ToBytes()
		path := "/hi"
		bucket := "paths"
		database.AddElementToBucket(db, []byte(bucket), []byte(path), serialised)
		fmt.Println(count, "("+cloudFile.GoogleId+")")
	}
	fmt.Println("All messages printed, total received: ", count)
}
