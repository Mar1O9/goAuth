package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maro1O9/goauth/internal/database"
	"github.com/Maro1O9/goauth/internal/database/models"
	"github.com/Maro1O9/goauth/internal/handlers"
	"github.com/gin-gonic/gin"
)

func TestSignUp(t *testing.T) {
	database.MakeDb(&models.User{})
	tests := []struct {
		name       string
		input      handlers.InputUser
		statusCode int
	}{
		{
			name: "invalid JSON input",
			input: handlers.InputUser{
				Username:        "test",
				Name:            "Test",
				Email:           "test@example.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid username",
			input: handlers.InputUser{
				Username:        "i",
				Name:            "Test",
				Email:           "test@example.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid name",
			input: handlers.InputUser{
				Username:        "test",
				Name:            "in",
				Email:           "test@example.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			input: handlers.InputUser{
				Username:        "test",
				Name:            "Test",
				Email:           "invalid",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			input: handlers.InputUser{
				Username:        "test",
				Name:            "Test",
				Email:           "test@example.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "mismatched passwords",
			input: handlers.InputUser{
				Username:        "test",
				Name:            "Test",
				Email:           "test@example.com",
				Password:        "password",
				ConfirmPassword: "invalid",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "successful sign up",
			input: handlers.InputUser{
				Username:        "test",
				Name:            "Test",
				Email:           "test@example.com",
				Password:        "Password123",
				ConfirmPassword: "Password123",
			},
			statusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonInput, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/sign-up", bytes.NewBuffer(jsonInput))
			handlers.SignUp(c)
			if w.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d, response body: %s", tt.statusCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestLogin(t *testing.T) {
	database.MakeDb(&models.User{})

	tests := []struct {
		name       string
		input      handlers.LoginUser
		statusCode int
	}{
		{
			name: "invalid JSON input",
			input: handlers.LoginUser{
				Email:    "invalid",
				Password: "password",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			input: handlers.LoginUser{
				Email:    "invalid_email",
				Password: "Password123",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			input: handlers.LoginUser{
				Email:    "test@example.com",
				Password: "invalid_password",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "non-existent user",
			input: handlers.LoginUser{
				Email:    "non_existent@example.com",
				Password: "Password123",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "incorrect password",
			input: handlers.LoginUser{
				Email:    "test@example.com",
				Password: "Incorrectpassword123",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "successful login",
			input: handlers.LoginUser{
				Email:    "test@example.com",
				Password: "Password123",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			mjsn, err := json.Marshal(test.input)
			if err != nil {
				t.Fatal(err)
			}
			c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(mjsn))

			handlers.Login(c)

			if w.Code != test.statusCode {
				t.Errorf("expected status code %d, got %d", test.statusCode, w.Code)
			}
		})
	}
}
