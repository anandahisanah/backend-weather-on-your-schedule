package forecast

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type responseFindForecastByDatetime struct {
	ID          int    `json:"id"`
	CityName    string `json:"city_name"`
	Datetime    string `json:"datetime"`
	Weather     string `json:"weather"`
	Humidity    string `json:"humidity"`
	WindSpeed   string `json:"wind_speed"`
	Temperature string `json:"temperature"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type responseGet struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Cities []city `json:"cities"`
}

type city struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func FindForecastByDatetime(c *gin.Context) {
	db := database.GetDB()

	paramUserUsername := c.Query("user_username")
	paramDate := c.Query("date")
	paramTime := c.Query("time")

	// Parse tanggal dengan format "dd/MM/yyyy"
	dateFormat := "02/01/2006"
	date, err := time.Parse(dateFormat, paramDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid date format",
			"original_message": err.Error(),
			"data":             nil,
		})
		return
	}

	// Parse waktu dengan format "15:04"
	timeFormat := "15:04"
	parsedTime, err := time.Parse(timeFormat, paramTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid time format",
			"original_message": err.Error(),
			"data":             nil,
		})
		return
	}

	// Mendapatkan lokasi waktu saat ini
	location, _ := time.LoadLocation("Local")

	// Menggabungkan tanggal dan waktu menjadi objek time.Time
	datetime := time.Date(date.Year(), date.Month(), date.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), location)

	// find user
	var user models.User
	if err := db.Where("username = ?", paramUserUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to get User",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// Check if data exists for the given date
	var count int64
	if err := db.Model(&models.Forecast{}).Where("city_id = ? AND DATE(datetime) = ?", user.CityID, date.Format("2006-01-02")).Count(&count).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to check data",
			"original_message": err.Error(),
			"data":             nil,
		})
		return
	}

	// If data not found for the given date, send response and stop further processing
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Data not found",
			"original_message": nil,
			"data":             nil,
		})
		return
	}

	// find forecast
	var forecast models.Forecast
	if err := db.Where("city_id = ? AND datetime = ?", user.CityID, datetime).Preload("City").First(&forecast).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":             400,
				"status":           "failed",
				"message":          "Forecast not found",
				"original_message": err.Error(),
				"data":             nil,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to get Forecast",
			"original_message": err.Error(),
			"data":             nil,
		})
		return
	}

	response := responseFindForecastByDatetime{
		ID:          int(forecast.ID),
		CityName:    forecast.City.Name,
		Datetime:    forecast.Datetime.String(),
		Weather:     forecast.Weather,
		Humidity:    forecast.Humidity,
		WindSpeed:   forecast.WindSpeed,
		Temperature: forecast.Temperature,
		CreatedAt:   forecast.CreatedAt.String(),
		UpdatedAt:   forecast.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    response,
	})
}


func GetForecast(c *gin.Context) {
	db := database.GetDB()

	paramCityName := c.Query("city_name")

	// find with cities
	var city models.City
	if err := db.Where("name = ?", paramCityName).First(&city).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to find City",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// get forecast
	var forecasts []models.Forecast
	if err := db.Where("city_id = ? AND datetime = ", city.ID).Find(&forecasts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to get Forecast",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// TODO: fix response
	// response
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    forecasts,
	})
}
