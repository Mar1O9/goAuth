package database_test

import (
	"testing"

	"github.com/Maro1O9/goauth/internal/database"
	"github.com/Maro1O9/goauth/internal/database/models"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	database.MakeDb(&models.User{})

	user := &models.User{
		Username:     "test",
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: []byte("Password123"),
	}
	err := database.Create(&models.User{}, user)
	assert.NoError(t, err)
	// Test Create with invalid table
	err = database.Create("invalid_table", user)
	assert.Error(t, err)
	// Test Create with invalid object
	err = database.Create(&models.User{}, "invalid_object")
	assert.Error(t, err)
	// Test Create with nil object
	err = database.Create(&models.User{}, nil)
	assert.Error(t, err)
	// Test Create with error from database
	database.DB.Exec("DROP TABLE users")
	err = database.Create("users", user)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	// Initialize the database
	database.MakeDb(&models.User{})
	// Test Update with valid model, object, query, and arguments
	user := &models.User{
		Username:     "test",
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: []byte("Password123"),
	}
	err := database.Create(&models.User{}, user)
	assert.NoError(t, err)
	err = database.Update(&models.User{}, &user, "id = ?", user.ID)
	assert.NoError(t, err)
	// Test Update with invalid model
	err = database.Update("invalid_model", &user, "id = ?", user.ID)
	assert.Error(t, err)
	// Test Update with invalid object
	err = database.Update(&models.User{}, "invalid_object", "id = ?", user.ID)
	assert.Error(t, err)
	// Test Update with invalid query
	err = database.Update(&models.User{}, &user, "invalid_query", user.ID)
	assert.Error(t, err)
	// Test Update with invalid arguments
	err = database.Update(&models.User{}, &user, "id = ?", "invalid_args")
	assert.Error(t, err)
	// Test Update with nil model
	err = database.Update(nil, &user, "id = ?", user.ID)
	assert.Error(t, err)
	// Test Update with nil object
	err = database.Update(&models.User{}, nil, "id = ?", user.ID)
	assert.Error(t, err)
	// Test Update with nil query
	err = database.Update(&models.User{}, &user, nil, user.ID)
	assert.Error(t, err)
	// Test Update with nil arguments
	err = database.Update(&models.User{}, &user, "id = ?", nil)
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	// Initialize the database
	database.MakeDb(&models.User{})

	// Test Delete with valid model, object, query, and arguments
	user := &models.User{
		Username:     "test",
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: []byte("Password123"),
	}
	err := database.Create(&models.User{}, user)
	assert.NoError(t, err)
	err = database.Delete(&models.User{}, user, "id = ?", user.ID)
	assert.NoError(t, err)

	// Test Delete with nil object
	err = database.Delete(&models.User{}, nil, "id = ?", 1)
	assert.Error(t, err)

	// Test Delete with nil database connection
	database.DB = nil
	err = database.Delete(&models.User{}, user, "id = ?", user.ID)
	assert.Error(t, err)
	database.MakeDb(&models.User{})

	// Test Delete with nil model
	err = database.Delete(nil, user, "id = ?", user.ID)
	assert.Error(t, err)

	// Test Delete with nil query
	err = database.Delete(&models.User{}, user, nil, 1)
	assert.Error(t, err)

	// Test Delete with invalid query
	err = database.Delete(&models.User{}, user, "invalid_query", 1)
	assert.Error(t, err)

	// Test Delete with no rows affected
	err = database.Delete(&models.User{}, user, "id = ?", 999)
	assert.Error(t, err)
}
