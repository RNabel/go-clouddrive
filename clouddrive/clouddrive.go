package clouddrive

import (
	"github.com/boltdb/bolt"
	"google.golang.org/api/drive/v3"
)

type CDrive struct {
	db *bolt.DB
	drv *drive.Service
}

func (cd *CDrive) GetRootId() string {
	f, _ := cd.drv.Files.Get("root").Do()
	return f.Id
}

func (cd *CDrive) DB() *bolt.DB {
	return cd.db
}

func (cd *CDrive) Drive() *drive.Service {
	return cd.drv
}

func NewCDrive(drv *drive.Service, db *bolt.DB) *CDrive {
	cd := new(CDrive)
	cd.drv = drv
	cd.db = db

	return cd
}
