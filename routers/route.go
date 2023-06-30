package routers

import (
	"os"
	"see-weather-on-your-schedule/controllers/event"
	"see-weather-on-your-schedule/controllers/province"
	"see-weather-on-your-schedule/controllers/user"

	"github.com/gin-gonic/gin"
)

var PORT = os.Getenv("PORT")

func StartServer() *gin.Engine {
	router := gin.Default()

	// user
	router.POST("/user", user.CreateUser)
	router.POST("/user/:id", user.UpdateUser)

	// province
	router.GET("/province", province.Find)

	// event
	router.POST("/event", event.CreateEvent)
	router.POST("/event/:id", event.UpdateEvent)

	// port
	if PORT == "" {
		PORT = "8080"
	}

	router.Run(":" + PORT)

	return router
}
