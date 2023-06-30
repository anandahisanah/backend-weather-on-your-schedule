package city

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseFind struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Find(c *gin.Context) {
	db := database.GetDB()

	// param
	paramName := c.Query("name")

	// find province
	var province models.Province
	if err := db.Where("name = ?", paramName).Find(&province).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to find Province",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// find city
	var cities []models.City
	if err := db.Where("province_id = ?", province.ID).Find(&cities).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to find City",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// looping
	var responses []responseFind
	for _, city := range cities {
		response := responseFind{
			ID:   int(city.ID),
			Name: city.Name,
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
