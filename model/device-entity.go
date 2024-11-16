package model

import (
	"github.com/google/uuid"
	"time"
)

type Device struct {
	ID        uint      `json:"-" gorm:"column:device_id;primarykey"`
	UUID      uuid.UUID `json:"uuid" gorm:"column:device_uuid;type:uuid;not null;"`
	Name      string    `json:"name" binding:"required" gorm:"column:device_name;type:text;not null;unique"`
	BrandName string    `json:"brand" binding:"required" gorm:"column:brand_name;type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"colum:created_at"`
}

func (b *Device) TableName() string {
	return "devices"
}
