package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	for {
		response, err := http.Get("http://localhost:8080/limit")
		if err != nil {
			log.Fatal(err)
		} else {
			body, _ := io.ReadAll(response.Body)
			fmt.Printf("status: %d,message: %s\n", response.StatusCode, body)
		}
		time.Sleep(time.Millisecond * 100 * time.Duration(rand.Intn(2)))
		//time.Sleep(time.Millisecond * 100)
	}
}
