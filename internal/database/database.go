package database

// database.go provides a singleton pattern for gorm.DB
// It is a wrapper around the gorm.DB object and provides
// methods for creating, updating and deleting records in the database

import (
	"errors"
	"log"
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var once sync.Once

var DB *gorm.DB

// MakeDb creates a database instance and migrates the given schema
// to the database. It ensures that the database instance is only
// created once. If the database instance already exists, it will not
// create a new one.
func MakeDb(schema ...interface{}) {
	if DB == nil {
		once.Do(
			func() {
				log.Println("Creating database instance now.")
				db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
				if err != nil {
					log.Fatalln(err.Error())
				}
				db.AutoMigrate(schema...)
				DB = db
			})
	} else {
		log.Println("Database instance already created.")
	}
}

// Create creates a new record in the database.
// It takes a model and an object and inserts the object into the database table
// that corresponds to the model. Returns an error if the create operation fails.
func Create(model interface{}, obj interface{}) error {
	result := DB.Model(model).Create(obj)
	return result.Error
}

// Update updates existing records in the database that match the provided query.
// It takes a model, an object with the updated fields, a query, and arguments for the query.
// Returns an error if the update operation fails.
func Update(model interface{}, obj, query interface{}, args ...interface{}) error {
	if DB == nil {
		return errors.New("database connection not established")
	}
	if model == nil {
		return errors.New("model cannot be nil")
	}
	if obj == nil {
		return errors.New("object cannot be nil")
	}
	if query == nil {
		return errors.New("query cannot be nil")
	}
	result := DB.Model(model).Where(query, args).Updates(obj)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// Delete deletes a record from the database that matches the provided object.
// It takes an object that must have a valid primary key set.
// Returns an error if the delete operation fails.
func Delete(model interface{}, obj interface{}, query interface{}, args ...interface{}) error {
	if obj == nil {
		return errors.New("object cannot be nil")
	}

	if DB == nil {
		return errors.New("database connection not established")
	}

	if model == nil {
		return errors.New("model cannot be nil")
	}

	if query == nil {
		return errors.New("query cannot be nil")
	}

	result := DB.Model(model).Where(query, args...).Delete(obj)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
