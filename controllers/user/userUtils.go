package user

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/simhonchourasia/betfr-be/authentication"
	"github.com/simhonchourasia/betfr-be/models"
)

func validateRegistration(registration models.Registration) error {
	if len(registration.Username) < 1 || len(registration.Username) > 30 {
		return fmt.Errorf("username must be between 1 and 30 characters")
	}
	if strings.Contains(registration.Username, "@") {
		return fmt.Errorf("username cannot contain the @ character")
	}
	if _, err := mail.ParseAddress(registration.Email); err != nil {
		return fmt.Errorf("must use valid email")
	}
	if len(registration.Password) < 6 || 100 < len(registration.Password) {
		return fmt.Errorf("password must be between 6 and 100 characters")
	}
	return nil
}

func createUserFromRegistration(registration models.Registration) (models.User, error) {
	var user models.User

	user.Email = registration.Email
	user.Username = registration.Username
	user.PasswordHash = authentication.HashPassword(registration.Password)
	token, refreshToken, err := authentication.GenerateAllTokens(user.Username)
	if err != nil {
		return user, err
	}
	user.Token = token
	user.RefreshToken = refreshToken

	user.RegistrationTime = time.Now()
	user.LastLoginTime = time.Now()
	// very funny
	user.ProfilePicLink = "https://cdn.discordapp.com/attachments/753412713920725015/1173437077594046534/blank-profile-picture-973460-1-1-1080x1080.png?ex=6563f370&is=65517e70&hm=41a65ee2e4e0696f3ba2e506f5ca60e14586f4ef07f6c0dd2b51c63ec91c43f8&"
	return user, nil
}
