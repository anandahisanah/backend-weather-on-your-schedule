package province

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type requestFirst struct {
	ID int `json:"id"`
}

type responseFind struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Cities []city `json:"cities"`
}

type city struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Find(c *gin.Context) {
	db := database.GetDB()

	// find with city
	var provinces []models.Province
	if err := db.Preload("City").Find(&provinces).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to get Province",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// looping
	var responses []responseFind
	for _, province := range provinces {
		var cities []city
		for _, cityLoop := range province.Cities {
			cityResponse := city{
				ID:   int(cityLoop.ID),
				Name: cityLoop.Name,
			}
			cities = append(cities, cityResponse)
		}
		response := responseFind{
			ID:     int(province.ID),
			Name:   province.Name,
			Cities: cities,
		}
		responses = append(responses, response)
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responses,
	})

}
