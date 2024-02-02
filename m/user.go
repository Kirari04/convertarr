package m

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	PwdHash   string
}
