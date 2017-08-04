package main

import (
	"sync"
	"os"
	"CloudDrive/cloudconn"
	"CloudDrive/database"
	"CloudDrive/files"
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

	if !database.IsDbInitialised(db) {
		// Copy all file metadata from GDrive to key value store.
		fileChan := make(chan files.CloudFile, 1000)
		wg.Add(1)
		go cloudconn.GetAllFilesFromDrive(drv, fileChan, wg)
		wg.Add(1)
		go files.AddFilesToDB(fileChan, db, wg)
	}

	wg.Wait()
}