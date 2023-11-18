package controllers

import "github.com/simhonchourasia/betfr-be/dbinterface"

// Allows dependency injection of DB

type Handler struct {
	Db dbinterface.DBInterface
}
