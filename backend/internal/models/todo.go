package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo is used for the PUBLIC (unauthenticated) JSON-file-backed store
type Todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

// UserTodo is used for the AUTHENTICATED per-user store (Postgres via GORM)
type UserTodo struct {
	ID        uint           `json:"id"        gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"         gorm:"index"`
	UserID    uint           `json:"user_id"   gorm:"not null;index"`
	Item      string         `json:"item"      gorm:"not null"`
	Completed bool           `json:"completed" gorm:"default:false"`
}
