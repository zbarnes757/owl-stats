package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // instantiates the postgres dialect
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	databaseURL := os.Getenv("database_url")

	conn, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	// Database migration
	db.AutoMigrate(&Account{}, &Polling{}, &OWLPlayer{})

	if GetPolling("player_fetch") == nil {
		playerFetch := Polling{Title: "player_fetch"}
		db.Create(&playerFetch)
	}
}

// GetDB retrieves a connection to the database
func GetDB() *gorm.DB {
	return db
}

// Base contains common columns for all tables.
type Base struct {
	CreatedAt time.Time  `jsonapi:"attr,createdAt"`
	UpdatedAt time.Time  `jsonapi:"attr,createdAt"`
	DeletedAt *time.Time `sql:"index" jsonapi:"attr,createdAt"`
}
