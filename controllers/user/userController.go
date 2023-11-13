package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/authentication"
	"github.com/simhonchourasia/betfr-be/config"
	"github.com/simhonchourasia/betfr-be/controllers"
	"github.com/simhonchourasia/betfr-be/models"
	"gorm.io/gorm"
)

type UserHandler controllers.Handler
type Yeah int

func (userHandler *UserHandler) getUserWithUsernameOrEmail(usernameOrEmail string) (models.User, error) {
	var user models.User
	var query string
	if strings.Contains(usernameOrEmail, "@") {
		// looking for email
		query = "email = ?"
	} else {
		// looking for username
		query = "username = ?"
	}

	err := userHandler.Db.Where(query, usernameOrEmail).First(&user)
	return user, err.Error
}

// Checks that the given username and email combination is unique
func (userHandler *UserHandler) isUserUnique(email string, username string) (bool, error) {
	var user models.User // unused
	err := userHandler.Db.Where("email = ?", email).Or("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (userHandler *UserHandler) SignUpFunc(c *gin.Context) {
	var registration models.Registration
	if err := c.BindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateRegistration(registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isUserUnique, err := userHandler.isUserUnique(registration.Email, registration.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isUserUnique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "An account with that username or email already exists"})
	}

	user, err := createUserFromRegistration(registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Something went wrong in user signup: %s", err.Error())})
	}

	if err := userHandler.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Something went wrong in user signup: %s", err.Error())})
	}

	c.JSON(http.StatusOK, "ok")
}

func (userHandler *UserHandler) LoginFunc(c *gin.Context) {
	var login models.Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userHandler.getUserWithUsernameOrEmail(login.UsernameOrEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordOk := authentication.VerifyPassword(login.Password, user.PasswordHash)
	if !passwordOk {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect password"})
	}

	token, refreshToken, err := authentication.GenerateAllTokens(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	authentication.UpdateAllTokens(userHandler.Db, token, refreshToken, user.ID)

	c.SetCookie("jwt", token, 24*60*60, "/", config.GlobalConfig.Domain, false, true)
	c.JSON(http.StatusOK, user)
}

// Get currently-signed in user from cookie
func (userHandler *UserHandler) GetUserFunc(c *gin.Context) {
	claims, statusCode, err := authentication.GetClaimsFromCookie(c)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	user, err := userHandler.getUserWithUsernameOrEmail(claims.Issuer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (userHandler *UserHandler) LogoutFunc(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", config.GlobalConfig.Domain, false, true)

	c.JSON(http.StatusOK, gin.H{"msg": "logged out"})
}
