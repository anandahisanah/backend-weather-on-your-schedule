package routers

import (
	"backend-weather-on-your-schedule/controllers/city"
	"backend-weather-on-your-schedule/controllers/event"
	"backend-weather-on-your-schedule/controllers/province"
	"backend-weather-on-your-schedule/controllers/user"
	"os"

	"github.com/gin-gonic/gin"
)

var PORT = os.Getenv("PORT")

func StartServer() *gin.Engine {
	router := gin.Default()

	// user
	router.GET("/user/login", user.Login)
	router.GET("/user/:id", user.FindUser)
	router.POST("/user", user.CreateUser)
	router.POST("/user/:id", user.UpdateUser)

	// province
	router.GET("/province", province.Find)

	// city
	router.GET("/city", city.Find)

	// event
	router.GET("/events", event.GetEvent)
	router.GET("/event/:id", event.FindEvent)
	router.POST("/event", event.CreateEvent)
	router.POST("/event/:id", event.UpdateEvent)

	// port
	if PORT == "" {
		PORT = "8080"
	}

	router.Run(":" + PORT)

	return router
}
