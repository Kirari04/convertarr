package m

import (
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type History struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Hash          string
	OldPath       string
	NewPath       string
	OldSize       uint64
	NewSize       uint64
	TimeTaken     time.Duration
	ComparisonImg string
	Progress      float64
	Error         string `gorm:"size:10000"`
	Status        string // encoding | failed | finished | copy
}

func (j *History) Create(DB *gorm.DB, OldPath string) error {
	j.OldPath = OldPath
	if err := DB.Create(j).Error; err != nil {
		log.Error("Failed to create history", err)
		return err
	}
	return nil
}

func (j *History) SetNewPath(DB *gorm.DB, NewPath string) error {
	j.NewPath = NewPath
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) SetProgress(DB *gorm.DB, Progress float64) error {
	j.Progress = Progress
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) SetComparisonImg(DB *gorm.DB, ComparisonImg string) error {
	j.ComparisonImg = ComparisonImg
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) SetOldSize(DB *gorm.DB, OldSize uint64) error {
	j.OldSize = OldSize
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) SetHash(DB *gorm.DB, Hash string) error {
	j.Hash = Hash
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) Encoding(DB *gorm.DB) error {
	j.Status = "encoding"
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) Failed(DB *gorm.DB, Error string) error {
	j.Error = Error
	j.Status = "failed"
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) Copy(DB *gorm.DB, NewPath string) error {
	j.NewPath = NewPath
	j.Status = "copy"
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}

func (j *History) Finished(DB *gorm.DB, NewSize uint64, TimeTaken time.Duration) error {
	j.NewSize = NewSize
	j.TimeTaken = TimeTaken
	j.Status = "finished"
	if err := DB.Save(j).Error; err != nil {
		log.Error("Failed to save history", err)
		return err
	}
	return nil
}
