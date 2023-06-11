package models

import "time"

type City struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int        `gorm:"not null" json:"province_id" valid:"required~province_id is required"`
	Province   Province   `gorm:"foreignKey:ProvinceID" json:"province"`
	Name       string     `gorm:"not null" json:"name" valid:"required~name is required"`
	Users      []User     `gorm:"foreignKey:ProvinceID" json:"users"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
