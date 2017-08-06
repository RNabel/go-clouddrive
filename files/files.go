package files

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"path"
	"sync"

	"CloudDrive/types"
	"bytes"
	"encoding/gob"
)

func AddFilesToDB(cd types.Drive, fileChan chan types.File, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get id of 'root'.
	rootId := cd.GetRootId()

	// Create map from ID to children. To be populated when there is no match
	//  in the previous map and
	idChildrenMap := make(map[string][]types.File)

	var count = 0

	// Build up dependency tree.
	for cloudFile := range fileChan {
		count += 1
		parents := cloudFile.Parents()

		for _, parentID := range parents {
			idChildrenMap[parentID] = append(idChildrenMap[parentID], cloudFile)
		}
	}

	fmt.Println("Downloads complete.")
	paths := getAllPathPairs(idChildrenMap, rootId, "/")

	bucket := "paths"
	for i, p := range paths {
		if i%1000 == 0 { // Print progress.
			fmt.Println("Added", i)
		}

		serialised := p.File.ToBytes()
		cd.DB().AddElementToBucket([]byte(bucket), []byte(p.Path), serialised)
	}
	fmt.Println("Database operations finished.")
}

type PathIdPair struct {
	File types.File
	Path string
}

func getAllPathPairs(idChildrenMap map[string][]types.File, rootID string, currentPath string) []PathIdPair {
	var ret = []PathIdPair{}

	for _, child := range idChildrenMap[rootID] {
		newPath := path.Join(currentPath, child.Name())
		newEntry := PathIdPair{child, newPath}
		ret = append(ret, newEntry)

		// Recursively add all other paths.
		other := getAllPathPairs(idChildrenMap, child.Id(), newPath)
		ret = append(ret, other...)
	}

	return ret
}

// Export types.
type CloudFile struct {
	googleId string
	name     string
	parents  []string
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

func (cf *CloudFile) CopyGoogleFile(file *drive.File) {
	cf.googleId = file.Id
	cf.name = file.Name
	cf.parents = file.Parents
}

func (cf *CloudFile) FromBytes(in []byte) {
	buf := bytes.Buffer{}
	buf.Write(in)

	decoder := gob.NewDecoder(&buf)
	err := decoder.Decode(cf)
	if err != nil {
		fmt.Println("failed to gob Decode", err)
	}
}

func (cf *CloudFile) Name() string {
	return cf.name
}

func (cf *CloudFile) Parents() []string {
	return cf.parents
}

func (cf *CloudFile) Id() string {
	return cf.googleId
}
