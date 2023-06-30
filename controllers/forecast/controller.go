package forecast

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseGet struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Cities []city `json:"cities"`
}

type city struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

	// response
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responses,
	})
}
