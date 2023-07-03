package forecast

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type responseFindForecastNowByCity struct {
	ID          int    `json:"id"`
	CityID      int    `json:"city_id"`
	CityName    string `json:"city_name"`
	Datetime    string `json:"datetime"`
	Weather     string `json:"weather"`
	Humidity    string `json:"humidity"`
	WindSpeed   string `json:"wind_speed"`
	Temperature string `json:"temperature"`
}

type responseGetForecast struct {
	ID          int    `json:"id"`
	CityID      int    `json:"city_id"`
	CityName    string `json:"city_name"`
	Datetime    string `json:"datetime"`
	Weather     string `json:"weather"`
	Humidity    string `json:"humidity"`
	WindSpeed   string `json:"wind_speed"`
	Temperature string `json:"temperature"`
}

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

func FindForecastNowByCity(c *gin.Context) {
	db := database.GetDB()

	paramUserUsername := c.Query("user_username")

	// find user
	var user models.User
	if err := db.Where("username = ?", paramUserUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to find User",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find forecast
	var forecast models.Forecast
	if err := db.Where("city_id = ? AND datetime >= ?", user.CityID, time.Now()).Preload("City").First(&forecast).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to find Forecast",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	response := responseFindForecastNowByCity{
		ID:          int(forecast.ID),
		CityID:      forecast.CityID,
		CityName:    forecast.City.Name,
		Datetime:    forecast.Datetime.String(),
		Weather:     forecast.Weather,
		Humidity:    forecast.Humidity,
		WindSpeed:   forecast.WindSpeed,
		Temperature: forecast.Temperature,
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

	paramCityId := c.Query("city_id")
	paramLimit := c.Query("limit")

	// get forecast
	var forecasts []models.Forecast
	var query *gorm.DB

	if paramLimit == "all" {
		query = db.Where("city_id = ? AND datetime >= ?", paramCityId, time.Now()).Preload("City").Find(&forecasts)
	} else {
		query = db.Where("city_id = ? AND datetime >= ?", paramCityId, time.Now()).Preload("City").Limit(5).Find(&forecasts)
	}

	if query.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to get Forecast",
			"original_message": query.Error.Error(),
			"data":             nil,
		})
		return
	}

	response := make([]responseFindForecastNowByCity, len(forecasts))
	for i, forecast := range forecasts {
		response[i] = responseFindForecastNowByCity{
			ID:          int(forecast.ID),
			CityID:      forecast.CityID,
			CityName:    forecast.City.Name,
			Datetime:    forecast.Datetime.String(),
			Weather:     forecast.Weather,
			Humidity:    forecast.Humidity,
			WindSpeed:   forecast.WindSpeed,
			Temperature: forecast.Temperature,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    response,
	})
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

	// find forecast
	var forecast models.Forecast
	if err := db.Where("city_id = ? AND datetime >= ?", user.CityID, datetime).Order("datetime ASC").Preload("City").First(&forecast).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to get Forecast",
			"original_message": err,
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
