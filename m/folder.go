package m

import (
	"time"
)

type Folder struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Path      string
}
