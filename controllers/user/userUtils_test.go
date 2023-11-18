package user

import (
	"testing"

	"github.com/simhonchourasia/betfr-be/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateRegistration(t *testing.T) {
	var testCases = []struct {
		name    string
		input   models.Registration
		isValid bool
	}{
		{"Standard case", models.Registration{Username: "simhon", Email: "s@c.com", Password: "qwertyuiop"}, true},
		{"Short password", models.Registration{Username: "asdf", Email: "s@c.com", Password: "no"}, false},
		{"Invalid email", models.Registration{Username: "asdf", Email: "asdf", Password: "no"}, false},
		{"Invalid username", models.Registration{Username: "as@df.co", Email: "asdf", Password: "no"}, false},
		{"Invalid username", models.Registration{Username: "", Email: "asdf", Password: "no"}, false},
		{"Invalid username", models.Registration{Username: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Email: "asdf", Password: "no"}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ans := validateRegistration(testCase.input)
			// check only for presence of or lack of error
			if (ans == nil) && !testCase.isValid {
				t.Errorf(
					"test case with username %s and email %s should not pass validation",
					testCase.input.Username,
					testCase.input.Email,
				)
			}
			if (ans != nil) && testCase.isValid {
				t.Errorf(
					"test case with username %s and email %s should pass validation",
					testCase.input.Username,
					testCase.input.Email,
				)
			}
		})
	}
}

func TestCreateUserFromRegistration(t *testing.T) {
	// TODO: add more checks here
	reg := models.Registration{Username: "raidennnn", Email: "s@c.com", Password: "qwertyuiop"}
	user, err := createUserFromRegistration(reg)
	assert.Nil(t, err)
	if user.PasswordHash == reg.Password {
		assert.NotEqual(t, user.PasswordHash, reg.Password, "Password hash should be different from password")
	}
}
