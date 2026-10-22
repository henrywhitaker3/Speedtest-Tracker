package database

import (
	"time"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `                            json:"created_at"`
	UpdatedAt time.Time      `                            json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                json:"deleted_at"`
}
