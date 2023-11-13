package controllers

import "gorm.io/gorm"

// Allows dependency injection of DB
type Handler struct {
	Db *gorm.DB
}
