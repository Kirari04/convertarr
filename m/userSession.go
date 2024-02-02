package m

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
	UserId    uuid.UUID `gorm:"index"`
}
