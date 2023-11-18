package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/config"
	"github.com/simhonchourasia/betfr-be/controllers"
	"github.com/simhonchourasia/betfr-be/controllers/friendship"
	"github.com/simhonchourasia/betfr-be/controllers/user"
	"github.com/simhonchourasia/betfr-be/middleware"
	"github.com/simhonchourasia/betfr-be/migrations"
	"github.com/simhonchourasia/betfr-be/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := config.SetupConfig()
	if err != nil {
		panic("Error in config: " + err.Error())
	}

	dsn := config.GlobalConfig.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// If migrationsOnly is set in the config JSON, we will only migrate and then exit
	if config.GlobalConfig.MigrationsOnly {
		migrations.MigrateUsers(db)
		return
	}

	router := gin.New()
	var handler controllers.Handler
	handler.Db = db

	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware)

	routes.UnprotectedUserRoutes(router, (user.UserHandler)(handler))
	routes.UnprotectedFriendshipRoutes(router, (friendship.FriendshipHandler)(handler))
	routes.UnprotectedBetRoutes(router)
	routes.UnprotectedStakeRoutes(router)

	// TESTING, REMOVE
	router.GET("/testing-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.Use(middleware.Authentication)

	// TESTING, REMOVE
	router.GET("/testing-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	routes.ProtectedUserRoutes(router, (user.UserHandler)(handler))
	routes.ProtectedFriendshipRoutes(router, (friendship.FriendshipHandler)(handler))
	routes.ProtectedBetRoutes(router)
	routes.ProtectedStakeRoutes(router)

	router.Run(":" + config.GlobalConfig.Port)

	fmt.Println("yeah")
}
