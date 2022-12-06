package routes

import (
	controller "golang-hotell-app/controllers"
	 "golang-hotell-app/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	 incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users" , controller.GetUsers())
	incomingRoutes.GET("users/:user_id", controller.GetUser())
}