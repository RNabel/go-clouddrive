package files

import (
	"google.golang.org/api/drive/v3"
	"fmt"
	"sync"
	"github.com/boltdb/bolt"
	"path"

	"CloudDrive/types"
	"CloudDrive/cloudconn"
)

func AddFilesToDB(fileChan chan types.CloudFile, db *bolt.DB, drv *drive.Service, wg sync.WaitGroup) {
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


	printAllPaths(idChildrenMap, rootId, "/")
	// TODO If not present add it to its parent's id.
	//serialised := cloudFile.ToBytes()
	//path := "/hi"
	//bucket := "paths"
	//fmt.Println(count, "("+cloudFile.GoogleId+")")
	//database.AddElementToBucket(db, []byte(bucket), []byte(path), serialised)
	// TODO Add all elements to the database
	//fmt.Println("All messages printed, total received: ", count)
}

type PathIdPair struct {
	File types.CloudFile
	Path string
}

func printAllPaths(idChildrenMap map[string][]types.CloudFile, rootID string, currentPath string) []PathIdPair {
	var ret = []PathIdPair{}

	for _, child := range idChildrenMap[rootID] {
		// FIXME child should be CloudFile element.
		newPath := path.Join(currentPath, child.Name)
		newEntry := PathIdPair{child, newPath}
		ret = append(ret, newEntry)
		fmt.Println(newEntry.Path, "(" + newEntry.File.GoogleId + ")")

		// Recursively add all other paths.
		other := printAllPaths(idChildrenMap, child.GoogleId, newPath)
		ret = append(ret, other...)
	}

	return ret
}
