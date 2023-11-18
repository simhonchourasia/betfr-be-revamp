package authentication

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/simhonchourasia/betfr-be/config"
	"github.com/simhonchourasia/betfr-be/dbinterface"
	"github.com/simhonchourasia/betfr-be/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Used to hold JWT info
type SignedDetails struct {
	Username string
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GlobalConfig.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}

func GenerateAllTokens(username string) (string, string, error) {
	expiryHours := 24
	if config.GlobalConfig.Debug {
		expiryHours = 168
	}
	log.Printf("Created JWT token expiring in %d hours\n", expiryHours)
	claims := &SignedDetails{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Issuer:    username,
			ExpiresAt: time.Now().Local().Add(time.Duration(expiryHours) * time.Hour).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			Issuer:    username,
			ExpiresAt: time.Now().Local().Add(time.Duration(24) * time.Duration(7) * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.GlobalConfig.SecretKey))
	refreshToken, err2 := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(config.GlobalConfig.SecretKey))
	log.Printf("Generated JWT tokens for %s", username)
	if err != nil || err2 != nil {
		err = fmt.Errorf("%v %v", err, err2)
		log.Panic(err)
		return "", "", err
	}

	return token, refreshToken, err
}

func UpdateAllTokens(db dbinterface.DBInterface, signedToken string, signedRefreshToken string, userID uuid.UUID) error {
	// This is a hack because i don't know how to mock crdbgorm
	// so just don't try to call this with a mock interface
	dbImpl, ok := db.(*gorm.DB)
	if !ok {
		panic(fmt.Sprintf("Must use GormDB implementation for authentication... type is %s", reflect.TypeOf(db)))
	}
	return crdbgorm.ExecuteTx(context.Background(), dbImpl, nil,
		func(tx *gorm.DB) error {
			var user models.User
			db.First(&user, userID)
			user.Token = signedToken
			user.RefreshToken = signedRefreshToken
			return db.Save(&user).Error
		},
	)
}

func HashPassword(password string) string {
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}
	return string(pwdBytes)
}

func VerifyPassword(givenPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(givenPassword))
	return err == nil
}

func GetClaimsFromCookie(c *gin.Context) (*jwt.StandardClaims, int, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	token, err := jwt.ParseWithClaims(
		cookie,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GlobalConfig.SecretKey), nil
		},
	)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	return claims, http.StatusOK, nil
}

func CheckUserPermissions(c *gin.Context, username string) error {
	name, ok := c.Get("username")
	nameStr, isString := name.(string)
	if !ok || !isString {
		return fmt.Errorf("could not get username from context")
	}
	fmt.Printf("current user: %s\n", nameStr)
	if ok && isString {
		if username == nameStr {
			return nil
		}
	}
	return fmt.Errorf("current user '%s' does not have the required permissions (want '%s')", nameStr, username)
}
