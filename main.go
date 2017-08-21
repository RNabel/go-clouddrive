package main

import (
	"fmt"
	"go-clouddrive/clouddrive"
	"go-clouddrive/cmd"
	"go-clouddrive/files"
	"go-clouddrive/types"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	cd := clouddrive.NewCDrive()
	db := cd.DB()
	fmt.Println(db)
	if !cd.DB().IsInitialised() {
		// Copy all file metadata from GDrive to key value store.
		fileChan := make(chan types.File, 1000)
		wg.Add(1)
		go cd.FetchMetadataFromGoogleDrive(fileChan, &wg)
		wg.Add(1)
		go files.AddFilesToDB(cd, fileChan, &wg)
	}
	wg.Wait()

	wg.Add(1)
	// Start cmd.
	go cmd.Start(cd, &wg)
	wg.Wait()
}
