package models

import "time"

type Forecast struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID  int        `gorm:"not null" json:"province_id" valid:"required~province_id is required"`
	Province    Province   `gorm:"foreignKey:ProvinceID" json:"province"`
	CityID      int        `gorm:"not null" json:"city_id" valid:"required~city_id is required"`
	City        City       `gorm:"foreignKey:CityID" json:"city"`
	Datetime    *time.Time `gorm:"not null" json:"datetime" valid:"required~datetime is required"`
	Humidity    string     `gorm:"not null" json:"humidity" valid:"required~humidity is required"`
	WindSpeed   string     `gorm:"not null" json:"wind_speed" valid:"required~wind_speed is required"`
	Temperature string     `gorm:"not null" json:"temperature" valid:"required~temperature is required"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
