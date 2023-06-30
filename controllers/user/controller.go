package user

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type requestCreate struct {
	ProvinceID int    `json:"province_id"`
	CityID     int    `json:"city_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
}

type responseCreate struct {
	ID           int    `json:"id"`
	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name"`
	CityID       int    `json:"city_id"`
	CityName     string `json:"city_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
}

type requestUpdate struct {
	ProvinceID int    `json:"province_id"`
	CityID     int    `json:"city_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
}

type responseUpdate struct {
	ID           int    `json:"id"`
	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name"`
	CityID       int    `json:"city_id"`
	CityName     string `json:"city_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
}

func CreateUser(c *gin.Context) {
	db := database.GetDB()

	// validate request format
	var request requestCreate
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

	// check if username exists
	var existingUser models.User
	if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Username already exists",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// validate province
	var province models.Province
	if err := db.First(&province, request.ProvinceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid Province",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// validate city
	var city models.City
	if err := db.Where("province_id = ?", province.ID).First(&city, request.CityID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid City",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// define model
	user := models.User{
		ProvinceID: request.ProvinceID,
		CityID:     request.CityID,
		Username:   request.Username,
		Password:   request.Password,
		Name:       request.Name,
	}

	// create
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Failed to create user",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	data := responseCreate{
		ID:           int(user.ID),
		ProvinceID:   user.ProvinceID,
		ProvinceName: province.Name,
		CityID:       user.CityID,
		CityName:     city.Name,
		Username:     user.Username,
		Password:     user.Password,
		Name:         user.Name,
	}

	// response
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"status":  "success",
		"message": "Success",
		"data":    data,
	})
}

func UpdateUser(c *gin.Context) {
	db := database.GetDB()

	// param
	paramID := c.Param("id")

	// validate request format
	var request requestUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid request body",
			"original_message": err,
		})
		return
	}

	// find user
	var user models.User
	if err := db.Preload("Province").Preload("City").First(&user, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "User not found",
			"original_message": err,
		})
		return
	}

	// check if username exists
	var existingUser models.User
	if err := db.Where("username = ? AND id != ?", request.Username, paramID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Username already exists",
			"original_message": err,
		})
		return
	}

	// validate province
	var province models.Province
	if err := db.First(&province, request.ProvinceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid Province",
			"original_message": err,
		})
		return
	}

	// validate city
	var city models.City
	if err := db.Where("province_id = ?", province.ID).First(&city, request.CityID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":             400,
			"status":           "failed",
			"message":          "Invalid City",
			"original_message": err,
		})
		return
	}

	// define model
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
			"code":             400,
			"status":           "failed",
			"message":          "Failed to update user",
			"original_message": err,
		})
		return
	}

	data := responseCreate{
		ID:           int(user.ID),
		ProvinceID:   user.ProvinceID,
		ProvinceName: province.Name,
		CityID:       user.CityID,
		CityName:     city.Name,
		Username:     user.Username,
		Password:     user.Password,
		Name:         user.Name,
	}

	// response
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"status":  "success",
		"message": "Success",
		"data":    data,
	})
}
