package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/sundayonah/go-jwt-project/controllers"
)

func AuthRoutes(incomingRoutes, _ gin.Engine) {
	incomingRoutes.POST("users/signup", controllers.Signup())
	incomingRoutes.POST("users/login", controllers.Login())
}

