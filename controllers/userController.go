package controllers

import (
	"net/http"
	"see-weather-on-your-schedule/database"
	"see-weather-on-your-schedule/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	db := database.GetDB()

	var user models.User
	// bind JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request body",
		})
		return
	}

	// check if username exists
	var existingUser models.User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		// Username sudah terdaftar
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Username already exists",
		})
		return
	}

	// create user
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Failed to create user",
		})
		return
	}

	// response
	c.JSON(http.StatusCreated, gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"password":    user.Password,
		"name":        user.Name,
		"province_id": user.ProvinceID,
		"city_id":     user.CityID,
	})
}

func UpdateUser(c *gin.Context) {
	db := database.GetDB()

	// param
	userID := c.Param("id")

	// find user
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  404,
			"error": "User not found",
		})
		return
	}

	// bind json
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request body",
		})
		return
	}

	// check if username exists
	var existingUser models.User
	if err := db.Where("username = ? AND id != ?", user.Username, userID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Username already exists",
		})
		return
	}

	// save
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Failed to update user",
		})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"password":    user.Password,
		"name":        user.Name,
		"province_id": user.ProvinceID,
		"city_id":     user.CityID,
	})
}
