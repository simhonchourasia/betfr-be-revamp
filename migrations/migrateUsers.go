package migrations

import (
	"github.com/simhonchourasia/betfr-be/models"
	"gorm.io/gorm"
)

// TODO: split this up
func MigrateUsers(db *gorm.DB) {
	// db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Friendship{})
}
