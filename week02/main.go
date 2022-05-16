package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-camp/week02/service"
	"github.com/pkg/errors"
)

func main() {
	http := Router()

	http.Run(":8999")
}

//curl localhost:8999/user/1
func Router() *gin.Engine {
	router := gin.Default()
	router.GET("/user/:id", GetUserHandler)
	return router
}

func GetUserHandler(c *gin.Context) {
	//入参验证
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, "参数错误")
		return
	}
	//获取数据
	userService := service.NewUser()
	user, err := userService.GetUser(ID)
	if err != nil {
		//处理记录日志...
		log.Printf("original error: %T,v: %v\n", errors.Cause(err), errors.Cause(err)) //root error
		log.Printf("stack trace:\n %+v \n", err)

		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
