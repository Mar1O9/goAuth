// Copyright 2025 Mahmoud Abdelrahman <mahmud.yousif.04@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


package utils

import (
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/Maro1O9/goauth/internal/inputs"
	"github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateSignupData compins all the input validation
// in on function
func ValidateSignupData(user *inputs.InputUser) error {

	// Validate username
	if err := ValidateUsername(user.Username); err != nil {
		return err
	}

	// Validate name
	if err := ValidateName(user.Name); err != nil {
		return err
	}

	// Validate email
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	// Validate password
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	// Validate confirm password
	if err := ValidatePassword(user.ConfirmPassword); err != nil {
		return err
	}

	// Check if passwords match
	if user.Password != user.ConfirmPassword {
		return errors.New("Passwords do not match")
	}
	return nil
}

// ValidateLoginData compins all the input validation
// in on function
func ValidateLoginData(user *inputs.LoginUser) error {

	// Validate email format
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	// Validate password format
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}
	return nil
}

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
