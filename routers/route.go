package routers

import (
	"see-weather-on-your-schedule/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/user", controllers.CreateUser)
	router.POST("/user/:id", controllers.UpdateUser)

	return router
}
