package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sundayonah/go-jwt-project/controllers"
)

func AuthRoutes(incomingRoutes, _ gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
