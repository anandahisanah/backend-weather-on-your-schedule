package routers

import (
	"see-weather-on-your-schedule/controllers/user"
	"see-weather-on-your-schedule/controllers/event"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	// user
	router.POST("/user", user.CreateUser)
	router.POST("/user/:id", user.UpdateUser)

	// event
	router.POST("/event", event.CreateEvent)
	router.POST("/event/:id", event.UpdateEvent)

	return router
}
