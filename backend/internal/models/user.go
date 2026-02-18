package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	// Primary Key: Manual UUID handling
	UserID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// Audit Timestamps: GORM automatically handles these if names match
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Enables GORM soft deletes
	// Fields: citext and text
	Email        string `gorm:"type:citext;unique;not null"`
	PasswordHash string `gorm:"type:text;not null"`
}
