package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func SlidingWindowLimit(sliding *SlidingWindowLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if sliding.IsLimited() {
			c.AbortWithStatusJSON(400, map[string]string{"message": "请求数量超过限制"})
			return
		}
	}
}

func main() {
	r := gin.Default()
	//win 1s,slot win 100ms,限制最大3个请求
	sliding := NewSliding(time.Second, 100*time.Millisecond, 3)

	r.Use(SlidingWindowLimit(sliding))

	r.GET("/limit", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Fatal(r.Run())

}
