package db

import (
	"path"

	"github.com/Pauloo27/neptune/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect(dataFolder string) {
	db, err := gorm.Open(sqlite.Open(path.Join(dataFolder, "db.sqlite")), &gorm.Config{})
	utils.HandleError(err, "Cannot connect to db")
	Database = db

	Database.AutoMigrate(&Artist{})
	Database.AutoMigrate(&Album{})
}
