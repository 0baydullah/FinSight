package transaction

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null;index"`
	CategoryID      uint      `gorm:"not null;index"`
	Type            string    `gorm:"size:20;not null"`
	Amount          int64     `gorm:"not null"`
	Description     *string   `gorm:"size:500"`
	TransactionDate time.Time `gorm:"not null;index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
