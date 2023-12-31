package user

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseLogin struct {
	ID           int    `json:"id"`
	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name"`
	CityID       int    `json:"city_id"`
	CityName     string `json:"city_name"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type responseFind struct {
	ID           int    `json:"id"`
	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name"`
	CityID       int    `json:"city_id"`
	CityName     string `json:"city_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type requestCreate struct {
	ProvinceName string `json:"province_name"`
	CityName     string `json:"city_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
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

func Login(c *gin.Context) {
	db := database.GetDB()

	// param
	paramUsername := c.Query("username")
	paramPassword := c.Query("password")

	// find user
	var user models.User
	if err := db.Preload("Province").Preload("City").Where("username = ? AND password = ?", paramUsername, paramPassword).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":             404,
			"status":           "failed",
			"message":          "User not found",
			"original_message": err,
			"data":             nil,
		})
		return
	}

	// response
	responseLogin := responseLogin{
		ID:           int(user.ID),
		ProvinceID:   user.ProvinceID,
		ProvinceName: user.Province.Name,
		CityID:       user.CityID,
		CityName:     user.City.Name,
		Username:     user.Username,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responseLogin,
	})
}

func FindUser(c *gin.Context) {
	db := database.GetDB()

	paramID := c.Param("id")

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

	responseFind := responseFind{
		ID:           int(user.ID),
		ProvinceID:   user.ProvinceID,
		ProvinceName: user.Province.Name,
		CityID:       user.CityID,
		CityName:     user.City.Name,
		Username:     user.Username,
		Password:     user.Password,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Success",
		"data":    responseFind,
	})
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
	if err := db.Where("name = ?", request.ProvinceName).First(&province).Error; err != nil {
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
	if err := db.Where("province_id = ? AND name = ?", province.ID, request.CityName).First(&city).Error; err != nil {
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
		ProvinceID: int(province.ID),
		CityID:     int(city.ID),
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
