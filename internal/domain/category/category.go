package category

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Name      string `gorm:"size:100;not null"`
	IsSystem  bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
