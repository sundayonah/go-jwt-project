package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sundayonah/go-jwt-project/controllers"
	"github.com/sundayonah/go-jwt-project/middleware"
)

func UserRoutes(incomingRoutes, *gin.HandlerFunc) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/user/:user_id", controller.GetUser())

}
