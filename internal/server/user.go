package server

import (
	"errors"
	"net/http"

	"github.com/Maro1O9/goauth/internal/database"
	"github.com/Maro1O9/goauth/internal/database/models"
	"github.com/Maro1O9/goauth/internal/inputs"
	"github.com/Maro1O9/goauth/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SignUp handles the /signup route.
//
// It expects a JSON payload containing the fields:
// - username: string
// - name: string
// - email: string
// - password: string
// - confirm_password: string
//
// It validates the input fields, hashes the password and creates a new user
// in the database. If any of the input fields are invalid, it will return
// a 400 error with the specific error message. If the passwords do not match,
// it will return a 400 error with the message "passwords do not match".
//
// If the user is created successfully, it will return a 201 status code with
// the newly created user as JSON.
func (s *Server) SignUp(c *gin.Context) {
	var input inputs.InputUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the input
	if err := utils.ValidateSignupData(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Check if username already exists
	if err := database.DB.First(&models.User{}).Where("username = ?", input.Username).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Username Already Taken"})
		return
	}

	// Check if email already exists
	if err := database.DB.First(&models.User{}).Where("email = ?", input.Email).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Creating the user
	user := &models.User{
		Username:     input.Username,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
	}

	// Create the user
	if err := database.Create(&models.User{}, &user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Success": "Signup successful"})
}

// Login handles the /login route.
//
// It takes a JSON payload with an email and password. If the email and password
// are valid, it returns a 200 status code with a JSON response containing a
// success message and sets a cookie with a JWT token. If the email or password
// is invalid, it returns a 401 status code with a JSON response containing an
// error message.
func (s *Server) Login(c *gin.Context) {
	var input inputs.LoginUser

	// Bind JSON payload to LoginUser struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate the input
	if err := utils.ValidateLoginData(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var user models.User

	// Check if user with the given email exists
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	// Compare provided password with the stored password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := utils.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the JWT token as a cookie
	c.SetCookie("Authorization", token, 3600*24*7, "/", "localhost", true, false)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
