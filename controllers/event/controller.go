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

type requestCreate struct {
	UserID      int    `json:"user_id"`
	Datetime    string `json:"datetime"`
	Title       string `json:"title"`
	Description string `json:"description"`
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

	// define model
	event := models.Event{
		UserID:      request.UserID,
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
