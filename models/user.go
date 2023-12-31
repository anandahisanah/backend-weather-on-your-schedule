package models

import "time"

type User struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int        `gorm:"not null" json:"province_id" valid:"required~province_id is required"`
	Province   Province   `gorm:"foreignKey:ProvinceID" json:"province"`
	CityID     int        `gorm:"not null" json:"city_id" valid:"required~city_id is required"`
	City       City       `gorm:"foreignKey:CityID" json:"city"`
	Username   string     `gorm:"not null;unique" json:"username" valid:"required~username is required"`
	Password   string     `gorm:"not null" json:"password" valid:"required~password is required"`
	Name       string     `gorm:"not null" json:"name" valid:"required~name is required"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
