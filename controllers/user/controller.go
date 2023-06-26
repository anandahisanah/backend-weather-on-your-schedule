package user

import (
	"net/http"
	"see-weather-on-your-schedule/database"
	"see-weather-on-your-schedule/models"

	"github.com/gin-gonic/gin"
)

type requestCreate struct {
	ProvinceID int    `json:"province_id"`
	CityID     int    `json:"city_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
}

type requestUpdate struct {
	ProvinceID int    `json:"province_id"`
	CityID     int    `json:"city_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
}

func CreateUser(c *gin.Context) {
	db := database.GetDB()

	var userCreate requestCreate
	if err := c.BindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request payload",
		})
		return
	}

	// Assign relation Province and City before Create
	var province models.Province
	if err := db.First(&province, userCreate.ProvinceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid province ID",
		})
		return
	}

	var city models.City
	if err := db.First(&city, userCreate.CityID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid city ID",
		})
		return
	}

	user := models.User{
		ProvinceID: userCreate.ProvinceID,
		CityID:     userCreate.CityID,
		Username:   userCreate.Username,
		Password:   userCreate.Password,
		Name:       userCreate.Name,
	}

	// Create user
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
		"province_id": user.ProvinceID,
		"province":    province.Name,
		"city_id":     user.CityID,
		"city":        city.Name,
		"username":    user.Username,
		"password":    user.Password,
		"name":        user.Name,
	})
}

func UpdateUser(c *gin.Context) {
	db := database.GetDB()

	// param
	paramID := c.Param("id")

	// find user
	var user models.User
	if err := db.Preload("Province").Preload("City").First(&user, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  404,
			"error": "User not found",
		})
		return
	}

	// bind json
	var request requestUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid request body",
		})
		return
	}

	// check if username exists
	var existingUser models.User
	if err := db.Where("username = ? AND id != ?", request.Username, paramID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Username already exists",
		})
		return
	}

	// Assign relasi Province dan City sebelum Update
	var province models.Province
	if err := db.First(&province, user.ProvinceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid province ID",
		})
		return
	}

	var city models.City
	if err := db.First(&city, user.CityID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "Invalid city ID",
		})
		return
	}

	// Update user
	user.ProvinceID = city.ProvinceID
	user.CityID = int(city.ID)
	user.Username = request.Username
	user.Password = request.Password
	user.Name = request.Name
	user.Province = province
	user.City = city

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
		"province_id": user.ProvinceID,
		"province":    user.Province.Name,
		"city_id":     user.CityID,
		"city":        user.City.Name,
		"username":    user.Username,
		"password":    user.Password,
		"name":        user.Name,
	})
}
