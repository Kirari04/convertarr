package m

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	PwdHash   string
}
