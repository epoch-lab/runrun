package internal

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// DB is the global database connection pool.
var DB *gorm.DB

// User represents the data model for a user in the database.
// It corresponds to the `users` table.
type User struct {
	ID                uint      `gorm:"primaryKey"`
	Account           string    `gorm:"unique;not null"`
	Password          string    `gorm:"not null"` // 明文密码
	CurrentDistance   float64   `gorm:"default:0"`
	TargetDistance    float64   `gorm:"default:80.0"`
	IsRunningRequired bool      `gorm:"default:true"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

// InitDB initializes the database connection and auto-migrates the schema.
func InitDB() {
	var err error
	// Connect to the SQLite database.
	// It will create a file named `runrun.db` in the project root.
	DB, err = gorm.Open(sqlite.Open("runrun.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto-migrate the User schema.
	// This will create the `users` table if it doesn't exist.
	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("FATAL: Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Database schema migrated successfully.")
}