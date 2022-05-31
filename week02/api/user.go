package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-camp/week02/code"
	"github.com/go-camp/week02/service"
	"github.com/pkg/errors"
)

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
	if errors.Is(err, code.ErrNotFound) {
		//处理指定业务关注的错误
		//记录日志...
		log.Printf("original error: %T,v: %v\n", errors.Cause(err), errors.Cause(err)) //root error
		log.Printf("stack trace:\n %+v \n", err)

		c.JSON(http.StatusOK, gin.H{
			"code":    code.NotFound,
			"message": err.Error(),
		})
		return
	} else if err != nil {
		//处理其他错误
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    code.SystemErr,
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}
