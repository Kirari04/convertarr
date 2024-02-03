package m

import (
	"time"
)

type UserSession struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
	UserId    uint `gorm:"index"`
}
