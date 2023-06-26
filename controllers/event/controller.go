package event

import (
	"log"
	"net/http"
	"see-weather-on-your-schedule/database"
	"see-weather-on-your-schedule/models"
	"time"

	"github.com/gin-gonic/gin"
)

type requestCreate struct {
	UserID      int        `json:"user_id"`
	CityID      int        `json:"city_id"`
	Datetime    *time.Time `json:"datetime"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

type requestUpdate struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	CityID      int        `json:"city_id"`
	Datetime    *time.Time `json:"datetime"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

func CreateEvent(c *gin.Context) {
	db := database.GetDB()

	// validate request format
	var request requestCreate
	if err := c.BindJSON(&request); err != nil {
		// TODO: standard json response is code, status, message and data
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request payload",
		})
		return
	}

	// find forecast
	var forecast models.Forecast
	if err := db.Where("city_id == ? AND datetime == ?", request.CityID, request.Datetime).First(&forecast).Error; err != nil {
		log.Fatalln("Error find forecast:", err)
	}

	// define event
	event := models.Event{
		UserID:      request.UserID,
		ForecastID:  int(forecast.ID),
		Datetime:    request.Datetime,
		Title:       request.Title,
		Description: request.Description,
	}

	// create
	if err := db.Create(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "failed",
			"message": "Failed to create event",
		})
		return
	}

	// response
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"status":  "success",
		"message": "success",
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
			"code":  400,
			"error": "Invalid request payload",
		})
		return
	}

	// find forecast
	var forecast models.Forecast
	if err := db.Where("city_id == ? AND datetime == ?", request.CityID, request.Datetime).First(&forecast).Error; err != nil {
		log.Fatalln("Error find forecast:", err)
	}

	// find event
	var event models.Event
	if err := db.First(&event, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  404,
			"error": "Event not found",
		})
		return
	}

	// update event properties
	event.UserID = request.UserID
	event.ForecastID = int(forecast.ID)
	event.Datetime = request.Datetime
	event.Title = request.Title
	event.Description = request.Description

	// save
	if err := db.Save(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "failed",
			"message": "Failed to update event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"status": "success",
		"data":   event,
	})
}
