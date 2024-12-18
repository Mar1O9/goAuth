package utils

import (
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateUsername checks if the provided Username is valid.
// The Username must be between 3 and 32 characters long and may contain
// alphanumeric characters and the following special characters:
// !@#$%^&*()_+={[]}:;,.<>?/-
// Returns an error if the Username is invalid.
func ValidateUsername(Username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;,.<>?/-]{3,32}$`)

	if !re.MatchString(Username) {
		return errors.New("invalid Username")
	}
	return nil

}

// ValidateName checks if the provided Name is valid.
// The Name must be between 3 and 32 characters long and may contain
// alphanumeric characters and the following special characters:
// !@#$%^&*()_+={}\[\]:;,.<>?/-
// Returns an error if the Name is invalid.
func ValidateName(Name string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;,.<>?/-]{3,32}$`)
	if !re.MatchString(Name) {
		return errors.New("invalid Name")
	}
	return nil
}

// ValidateEmail checks if the provided email address is valid.
// The email must follow the pattern of a standard email address and must
// be between 3 and 32 characters long. Returns an error if the email is invalid.
func ValidateEmail(email string) error {
	eRe := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(eRe)
	if !re.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

// ValidatePassword checks if the password is valid. The password must be at
// least 8 characters long, must contain at least one uppercase letter and
// one number, and may contain the following characters: A-Z, a-z, 0-9, _, !,
// @, #, $, ^, &, *, (, ), and +. Returns an error if the password is
// invalid.
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	PRe := `^(?=.*[A-Z])(?=.*\d)[A-Za-z\d!@#$%^&*()_+]{8,32}$`
	re, err := regexp2.Compile(PRe, 0)
	if err != nil {
		return errors.New("failed to compile regex")
	}

	match, err := re.MatchString(password)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid password")
	}
	return nil
}

var SecretKey = []byte(os.Getenv("SECRET_KEY"))

// CreateToken creates a JWT token that is valid for 24 hours and contains the
// provided username. The token is signed with the secret key. Returns an error
// if the username is empty.
func CreateToken(email string) (string, error) {
	if email == "" {
		return "", errors.New("email cannot be zero")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sup": email,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString(SecretKey)
}

// VerifyToken takes a JWT token and verifies its validity. If the token is valid,
// it returns nil. If the token is invalid, it returns an error.
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
