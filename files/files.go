package files

import (
	"google.golang.org/api/drive/v3"
	"fmt"
	"sync"
	"github.com/boltdb/bolt"
	"path"

	"CloudDrive/types"
	"CloudDrive/cloudconn"
	"CloudDrive/database"
)

func AddFilesToDB(fileChan chan types.CloudFile, db *bolt.DB, drv *drive.Service, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get id of 'root'.
	rootId := cloudconn.GetRootId(drv)

	// Create map from ID to children. To be populated when there is no match
	//  in the previous map and
	idChildrenMap := make(map[string][]types.CloudFile)

	var count = 0

	// Build up dependency tree.
	for cloudFile := range fileChan {
		count += 1
		parents := cloudFile.Parents

		for _, parentID := range parents {
			idChildrenMap[parentID] = append(idChildrenMap[parentID], cloudFile)
		}
	}

	fmt.Println("Downloads complete.")
	paths := getAllPathPairs(idChildrenMap, rootId, "/")

	bucket := "paths"
	for i, p := range paths {
		if i % 1000 == 0 { // Print progress.
			fmt.Println("Added", i)
		}

		serialised := p.File.ToBytes()
		database.AddElementToBucket(db, []byte(bucket), []byte(p.Path), serialised)
	}
	fmt.Println("Database operations finished.")
}

type PathIdPair struct {
	File types.CloudFile
	Path string
}

func getAllPathPairs(idChildrenMap map[string][]types.CloudFile, rootID string, currentPath string) []PathIdPair {
	var ret = []PathIdPair{}

	for _, child := range idChildrenMap[rootID] {
		// FIXME child should be CloudFile element.
		newPath := path.Join(currentPath, child.Name)
		newEntry := PathIdPair{child, newPath}
		ret = append(ret, newEntry)

		// Recursively add all other paths.
		other := getAllPathPairs(idChildrenMap, child.GoogleId, newPath)
		ret = append(ret, other...)
	}

	return ret
}
