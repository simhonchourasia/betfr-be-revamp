package migrations

import (
	"github.com/simhonchourasia/betfr-be/models"
	"gorm.io/gorm"
)

func MigrateUsers(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
