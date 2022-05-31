package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-camp/week02/api"
)

func main() {
	http := Router()
	http.Run(":8999")
}

//curl localhost:8999/user/1
func Router() *gin.Engine {
	router := gin.Default()
	router.GET("/user/:id", api.GetUserHandler)
	return router
}
