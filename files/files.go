package files

import (
	"google.golang.org/api/drive/v3"
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
	"github.com/boltdb/bolt"
	"path"

	"CloudDrive/cloudconn"
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

func AddFilesToDB(fileChan chan CloudFile, db *bolt.DB, drv *drive.Service, wg sync.WaitGroup) {
	defer wg.Done()

	// Get id of 'root'.
	rootId := cloudconn.GetRootId(drv)

	// Create map from IDs to paths.
	idPathMap := make(map[string]string)
	idPathMap[rootId] = ""
	// Create map from ID to children. To be populated when there is no match
	//  in the previous map and
	idChildrenMap := make(map[string][]string)

	var count = 0

	for cloudFile := range fileChan {
		count += 1
		var parentDir = ""
		parents := cloudFile.Parents

		// TODO Check if the parent has been assigned a path.
		// TODO If it does, ensure all dependencies are resolved. (see idChildrenMap)
		// TODO 	AddElementToBucket when resolving path.
		if len(cloudFile.Parents) > 0 {
			parentDir = ""
		} else if mapParDir, ok := idPathMap[cloudFile.GoogleId]; ok {
			parentDir = mapParDir
		}

		idPathMap[cloudFile.GoogleId] = path.Join(parentDir, cloudFile.Name)

		// TODO If not present add it to its parent's id.

		serialised := cloudFile.ToBytes()
		path := "/hi"
		bucket := "paths"
		fmt.Println(count, "("+cloudFile.GoogleId+")")
		database.AddElementToBucket(db, []byte(bucket), []byte(path), serialised)
	}
	// TODO Add all elements to the database
	//fmt.Println("All messages printed, total received: ", count)
}
