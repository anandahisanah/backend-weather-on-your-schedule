package models

import "time"

type Province struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string     `gorm:"not null;unique" json:"code" valid:"required~code is required"`
	Name      string     `gorm:"not null;unique" json:"name" valid:"required~name is required"`
	Endpoint  string     `gorm:"not null;unique" json:"endpoint" valid:"required~endpoint is required"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
