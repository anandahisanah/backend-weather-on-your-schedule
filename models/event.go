package models

import "time"

type Event struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      User       `gorm:"not null" json:"user_id" valid:"required~user_id is required"`
	User        User       `gorm:"foreignKey:UserID" json:"user"`
	ForecastID  Forecast   `gorm:"not null" json:"forecast_id" valid:"required~forecast_id is required"`
	Forecast    Forecast   `gorm:"foreignKey:ForecastID" json:"forecast"`
	Datetime    *time.Time `gorm:"not null" json:"datetime" valid:"required~datetime is required"`
	Title       string     `gorm:"not null" json:"title" valid:"required~title is required"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
