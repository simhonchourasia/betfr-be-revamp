package migrations

import (
	"github.com/simhonchourasia/betfr-be/models"
	"gorm.io/gorm"
)

// TODO: split this up
func MigrateUsers(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}

func MigrateFriendships(db *gorm.DB) {
	db.AutoMigrate(&models.Friendship{})
}

func MigrateBets(db *gorm.DB) {
	db.AutoMigrate(&models.Bet{})
}

func MigrateStakes(db *gorm.DB) {
	db.AutoMigrate(&models.Stake{})
}
