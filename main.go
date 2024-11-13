package main

import (
	"fmt"
	"os"

	routes "github.com/sundayonah/go-jwt-project/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted for api-1",
		})
	})
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted for api-2",
		})
	})
	router.Run(":" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	fmt.Printf("Server started on port %s...\n", port)
}
