package backup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yusufwira/lern-golang-gin/controller"
	"github.com/yusufwira/lern-golang-gin/service"
)

var db = make(map[string]string)
var (
	userService    service.UserService       = service.New()
	UserController controller.UserController = controller.New(userService)
)

func UserRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/getUser", func(c *gin.Context) {
		c.JSON(http.StatusOK, UserController.FindAll())
	})

	r.POST("/postUser", func(c *gin.Context) {
		c.JSON(http.StatusOK, UserController.Save(c))
	})
	return r
}

func main() {
	r := UserRouter()
	r.Run(":8088")
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusNotFound, "pong")
	})

	r.GET("/getData", func(ctx *gin.Context) {
		var pakerjaan = make([]string, 2)
		pakerjaan[0] = "Magang"
		pakerjaan[1] = "Programer anak pak doni"
		ctx.JSON(http.StatusOK, gin.H{
			"muzadi":    "mantab",
			"umur":      "50",
			"pekerjaan": pakerjaan,
		})
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	return r
}
