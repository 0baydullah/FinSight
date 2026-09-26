package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Username  string `gorm:"size:50;not null;uniqueIndex"`
	Email     string `gorm:"size:150;not null;uniqueIndex"`
	Phone     string `gorm:"size:15;uniqueIndex"`
	Password  string `gorm:"size:255;not null"`
	Balance   int64  `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
