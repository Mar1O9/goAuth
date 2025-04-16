// Copyright 2025 Mahmoud Abdelrahman <deprecated>
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


package database_test

import (
	"os"
	"testing"

	"github.com/Maro1O9/goauth/internal/database"
	"github.com/Maro1O9/goauth/internal/database/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	database.MakeDb(&models.User{})

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate(t *testing.T) {

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
	// Test Update with valid model, object, query, and arguments
	user := &models.User{
		Username:     "test",
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: []byte("Password123"),
	}
	err := database.Create(&models.User{}, &user)
	require.NoError(t, err)

	tests := []struct {
		name    string
		model   interface{}
		object  interface{}
		query   string
		args    []interface{}
		wantErr bool
	}{
		{
			name:    "valid model, object, query, and arguments",
			model:   &models.User{},
			object:  &models.User{ID: 1, Name: "Test2"},
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: false,
		},
		{
			name:    "invalid model",
			model:   "invalid_model",
			object:  &user,
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "invalid object",
			model:   &models.User{},
			object:  "invalid_object",
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "invalid query",
			model:   &models.User{},
			object:  &user,
			query:   "invalid_query",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "invalid arguments",
			model:   &models.User{},
			object:  &user,
			query:   "id = ?",
			args:    []interface{}{"invalid_args"},
			wantErr: true,
		},
		{
			name:    "nil model",
			model:   nil,
			object:  &user,
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "nil object",
			model:   &models.User{},
			object:  nil,
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "nil query",
			model:   &models.User{},
			object:  &user,
			query:   "",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "nil arguments",
			model:   &models.User{},
			object:  &user,
			query:   "id = ?",
			args:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := database.Update(test.model, test.object, test.query, test.args...)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	database.MakeDb(&models.User{})
	user := &models.User{
		ID:           1,
		Username:     "test",
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: []byte("Password123"),
	}

	_ = database.Create(&models.User{}, user)

	tests := []struct {
		name      string
		model     interface{}
		obj       interface{}
		query     interface{}
		args      []interface{}
		wantErr   bool
		setupFunc func()
	}{
		{
			name:    "valid model, object, query, and arguments",
			model:   &models.User{},
			obj:     user,
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: false,
		},
		{
			name:    "nil object",
			model:   &models.User{},
			obj:     nil,
			query:   "id = ?",
			args:    []interface{}{1},
			wantErr: true,
		},
		{
			name:    "nil database connection",
			model:   &models.User{},
			obj:     user,
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: true,
			setupFunc: func() {
				database.DB = nil
			},
		},
		{
			name:    "nil model",
			model:   nil,
			obj:     user,
			query:   "id = ?",
			args:    []interface{}{user.ID},
			wantErr: true,
		},
		{
			name:    "nil query",
			model:   &models.User{},
			obj:     user,
			query:   nil,
			args:    []interface{}{1},
			wantErr: true,
		},
		{
			name:    "invalid query",
			model:   &models.User{},
			obj:     user,
			query:   "invalid_query",
			args:    []interface{}{1},
			wantErr: true,
		},
		{
			name:    "no rows affected",
			model:   &models.User{},
			obj:     user,
			query:   "id = ?",
			args:    []interface{}{999},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupFunc != nil {
				test.setupFunc()
				database.MakeDb(&models.User{})
			}
			err := database.Delete(test.model, test.obj, test.query, test.args...)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
