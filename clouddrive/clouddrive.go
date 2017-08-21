package clouddrive

import (
	"go-clouddrive/cloudconn"
	"go-clouddrive/database"
	"go-clouddrive/files"
	"go-clouddrive/types"
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"os"
	"sync"
)

type CDrive struct {
	db  *database.CloudDB
	drv *drive.Service
}

func (cd *CDrive) Init() {
	// Initialise drive object.
	drv := cloudconn.GetDrive()
	cd.drv = drv

	// Initialise database.
	cd.db = &database.CloudDB{}
	err := cd.db.OpenDB()
	if err != nil { // Close if database connection failed.
		fmt.Println(err)
		os.Exit(1)
	}
}

func (cd *CDrive) Teardown() {
	// Close database connection.
	cd.DB().CloseDB()
}

// Google Drive functions.
func (cd *CDrive) GetRootId() string {
	f, _ := cd.drv.Files.Get("root").Do()
	return f.Id
}

func (cd *CDrive) FetchMetadataFromGoogleDrive(output chan types.File, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initial request.
	var r, err = cd.Drive().Files.List().
		Fields("nextPageToken, files(id, name, parents)").
		PageSize(1000).
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	var counter = 1

	// Send all files to the passed to the output channel, going through all pages.
	var cont = true
	for cont {
		fmt.Println("Fetched page:", counter, "files in page:", len(r.Files))
		counter++
		for _, f := range r.Files {
			fout := files.CloudFile{}
			fout.CopyGoogleFile(f)
			output <- &fout
		}

		cont = r.NextPageToken != "" || len(r.Files) == 0

		// Request next page.
		r, err = cd.Drive().Files.
			List().
			Fields("nextPageToken, files(id, name, parents)").
			PageToken(r.NextPageToken).
			PageSize(1000).
			Do()
	}

	close(output)
}

// Getters.
func (cd *CDrive) DB() *database.CloudDB {
	return cd.db
}

func (cd *CDrive) Drive() *drive.Service {
	return cd.drv
}

// Constructor.
func NewCDrive() *CDrive {
	cd := new(CDrive)
	cd.Init()
	return cd
}
