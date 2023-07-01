package event

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type responseGet struct {
	ID          int    `json:"id"`
	Datetime    string `json:"datetime"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
	Title       string `json:"title"`
}

type responseFind struct {
	ID                  int    `json:"id"`
	UserID              int    `json:"user_id"`
	UserUsername        string `json:"user_username"`
	ForecastID          int    `json:"forecast_id"`
	ForecastWeather     string `json:"forecast_weather"`
	ForecastHumidity    string `json:"forecast_humidity"`
	ForecastWindSpeed   string `json:"forecast_wind_speed"`
	ForecastTemperature string `json:"forecast_temperature"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type requestCreate struct {
	UserUsername string    `json:"user_username"`
	Datetime     string `json:"datetime"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type responseCreate struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Datetime    string `json:"datetime"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type requestUpdate struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Datetime    string `json:"datetime"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type responseUpdate struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Datetime    string `json:"datetime"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func GetEvent(c *gin.Context) {
	db := database.GetDB()

	paramUserUsername := c.Param("userUsername")

	// find user
	var user models.User
	if err := db.Where("username = ?", paramUserUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "User not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find event
	var events []models.Event
	if err := db.Where("user_id = ?", user.ID).Preload("Forecast").Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "Event not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// response
	var responses []responseGet
	for _, event := range events {
		responseGet := responseGet{
			ID:          int(event.ID),
			Datetime:    event.Datetime.String(),
			Weather:     event.Forecast.Weather,
			Temperature: event.Forecast.Temperature,
			Title:       event.Title,
		}
		responses = append(responses, responseGet)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responses,
	})
}

func FindEvent(c *gin.Context) {
	db := database.GetDB()

	paramID := c.Param("id")

	// find event
	var event models.Event
	if err := db.Preload("User").Preload("Forecast").First(&event, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "Event not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// response
	responseFind := responseFind{
		ID:                  int(event.ID),
		UserID:              event.UserID,
		UserUsername:        event.User.Username,
		ForecastID:          event.ForecastID,
		ForecastWeather:     event.Forecast.Weather,
		ForecastHumidity:    event.Forecast.Humidity,
		ForecastWindSpeed:   event.Forecast.WindSpeed,
		ForecastTemperature: event.Forecast.Temperature,
		CreatedAt:           event.CreatedAt.String(),
		UpdatedAt:           event.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responseFind,
	})
}

func CreateEvent(c *gin.Context) {
	db := database.GetDB()

	// validate request format
	var request requestCreate
	if err := c.BindJSON(&request); err != nil {
		// TODO: standard json response is code, status, message, original_message and data
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid request body",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find user
	var user models.User
	if err := db.Where("username = ?", request.UserUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "User not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// change format datetime
	layout := "2006-01-02 15:04:05"
	loc := time.FixedZone("Asia/Singapore", 8*60*60) // Offset waktu UTC+8 (Waktu Standar Singapura)
	datetime, err := time.ParseInLocation(layout, request.Datetime, loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":             500,
			"status":           "failed",
			"message":          "Error changing datetime format",
			"original_message": err,
			"data":             nil,
		})
		return
	}
	
	// find forecast and looking nearest datetime
	var forecast models.Forecast
	if err := db.Where("city_id = ? AND datetime = ?", user.CityID, request.Datetime).First(&forecast).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// if not found, find the last datetime
			if err := db.Where("city_id = ? AND datetime < ?", user.CityID, request.Datetime).Order("datetime desc").First(&forecast).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":             400,
					"status":           "failed",
					"message":          "Forecast not found",
					"original_message": err,
					"data":             nil,
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":             400,
				"status":           "failed",
				"message":          "Error occurred while querying forecast",
				"original_message": err,
				"data":             nil,
			})
			return
		}
	}

	// define forecast id
	var forecastID *int
	if forecast.ID != 0 {
		forecastIDValue := int(forecast.ID)
		forecastID = &forecastIDValue
	}

	// define model
	event := models.Event{
		UserID:      int(user.ID),
		ForecastID:  *forecastID,
		Datetime:    &datetime,
		Title:       request.Title,
		Description: request.Description,
	}

	// create
	if err := db.Create(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "failed",
			"message": "Failed to create event",
			"data":    nil,
		})
		return
	}

	responseCreate := responseCreate{
		ID:          int(event.ID),
		UserID:      event.UserID,
		Datetime:    event.Datetime.String(),
		Title:       event.Title,
		Description: event.Description,
		CreatedAt:   event.CreatedAt.String(),
	}

	// response
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"status":  "success",
		"message": "Success",
		"data":    responseCreate,
	})
}

func UpdateEvent(c *gin.Context) {
	db := database.GetDB()

	// param
	paramID := c.Param("id")

	// validate request format
	var request requestUpdate
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid request body",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find event
	var event models.Event
	if err := db.First(&event, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "Event not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find user
	var user models.User
	if err := db.First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "User not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// change format datetime
	layout := "2006-01-02 15:04:05"
	loc, err := time.LoadLocation("Asia/Singapore") // Ganti dengan zona waktu yang sesuai
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":             500,
			"status":           "failed",
			"message":          "Error loading location",
			"original_message": err,
			"data":             nil,
		})
		return
	}
	datetime, err := time.ParseInLocation(layout, request.Datetime, loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":             500,
			"status":           "failed",
			"message":          "Error changing datetime format",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find forecast and looking nearest datetime
	var forecast models.Forecast
	if err := db.Where("city_id = ? AND datetime = ?", user.CityID, request.Datetime).First(&forecast).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// if not found, find the last datetime
			if err := db.Where("city_id = ? AND datetime < ?", user.CityID, request.Datetime).Order("datetime desc").First(&forecast).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":             400,
					"status":           "failed",
					"message":          "Forecast not found",
					"original_message": err,
					"data":             nil,
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":             400,
				"status":           "failed",
				"message":          "Error occurred while querying forecast",
				"original_message": err,
				"data":             nil,
			})
			return
		}
	}

	// define forecast id
	var forecastID *int
	if forecast.ID != 0 {
		forecastIDValue := int(forecast.ID)
		forecastID = &forecastIDValue
	}

	// update event properties
	event.UserID = request.UserID
	event.ForecastID = *forecastID
	event.Datetime = &datetime
	event.Title = request.Title
	event.Description = request.Description

	// save
	if err := db.Save(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to update event",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// response
	responseUpdate := responseUpdate{
		ID:          int(event.ID),
		UserID:      event.UserID,
		Datetime:    event.Datetime.String(),
		Title:       event.Title,
		Description: event.Description,
		CreatedAt:   event.CreatedAt.String(),
		UpdatedAt:   event.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responseUpdate,
	})
}
