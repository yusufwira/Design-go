package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"

	"github.com/yusufwira/lern-golang-gin/connection"
	"github.com/yusufwira/lern-golang-gin/controller"
)

var (
	UserController controller.UserController = controller.New()
)

func main() {
	r := gin.Default()

	connection.Middleware()

	auth := r.Group("/oauth2")
	{
		auth.GET("/token", ginserver.HandleTokenRequest)
	}

	api := r.Group("/api")
	{
		fmt.Println("masuk")
		api.Use(ginserver.HandleTokenVerify())
		fmt.Println("masuk2")
		api.GET("/test", func(c *gin.Context) {
			ti, exists := c.Get(ginserver.DefaultConfig.TokenKey)
			if exists {
				c.JSON(http.StatusOK, ti)
				return
			}
			c.String(http.StatusOK, "not found")
		})

		api.GET("/getUserOuath", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.Index())
		})

		api.POST("/postUser", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.Store(c))
		})

	}

	r.Run(":9096")
}
