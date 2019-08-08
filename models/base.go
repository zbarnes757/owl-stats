package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // instantiates the postgres dialect
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
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
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}
