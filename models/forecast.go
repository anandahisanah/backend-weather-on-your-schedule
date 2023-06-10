package models

import "time"

type Forecast struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int        `gorm:"not null" json:"province_id" valid:"required~province_id is required"`
	Province   Province   `gorm:"foreignKey:ProvinceID" json:"province"`
	CityID     int        `gorm:"not null" json:"city_id" valid:"required~city_id is required"`
	City       City       `gorm:"foreignKey:CityID" json:"city"`
	Datetime   *time.Time `gorm:"not null" json:"datetime" valid:"required~datetime is required"`
	Weather    string     `gorm:"not null" json:"weather" valid:"required~weather is required"`
	Raining    string     `gorm:"not null" json:"raining" valid:"required~raining is required"`
	Uv         string     `gorm:"not null" json:"uv" valid:"required~uv is required"`
	Wind       string     `gorm:"not null" json:"wind" valid:"required~wind is required"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
