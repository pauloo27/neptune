package db

import (
	"github.com/Pauloo27/my-tune/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect(dataFolder string) {
	db, err := gorm.Open(sqlite.Open(dataFolder+"/db.sqlite"), &gorm.Config{})
	utils.HandleError(err, "Cannot connect to db")
	Database = db

	Database.AutoMigrate(&Artist{})
	Database.AutoMigrate(&Album{})
}
