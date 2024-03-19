package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-web-template/internal/service"
)

type Api struct {
	userService service.User
}

// NewApiController returns a new instance of the Api struct.
func NewApiController(
	userService service.User,
) *Api {
	return &Api{
		userService: userService,
	}
}

func (a *Api) register(rg *gin.RouterGroup) {
	rg.GET("/hello", a.hello)
}

func (a *Api) hello(c *gin.Context) {
	defer func() {
		log.Println("query done")
	}()

	user, err := a.userService.GetByName("gopher")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"username": user.Username,
		"profile":  user.Profile,
	})
}
