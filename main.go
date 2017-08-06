package main

import (
	"sync"
	"os"
	"CloudDrive/cloudconn"
	"CloudDrive/database"
	"CloudDrive/files"
	"CloudDrive/types"
	"CloudDrive/cmd"
	"CloudDrive/clouddrive"
)

func main() {
	wg := sync.WaitGroup{}

	// Initialise drive object.
	drv := cloudconn.GetDrive()

	// Initialise database.
	db, err := database.OpenDB()
	if err != nil {
		os.Exit(1)
	}
	defer database.CloseDB(db)

	cd := clouddrive.NewCDrive(drv, db)

	if !database.IsDbInitialised(db) {
		// Copy all file metadata from GDrive to key value store.
		fileChan := make(chan types.File, 1000)
		wg.Add(1)
		go cloudconn.GetAllFilesFromDrive(cd, fileChan, &wg)
		wg.Add(1)
		go files.AddFilesToDB(cd, fileChan, &wg)
	}
	wg.Wait()

	// Start cmd.
	go cmd.Start()
	wg.Wait()
}