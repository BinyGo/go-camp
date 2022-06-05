package server

import (
	"context"
	"net/http"

	"github.com/blog/internal/service"
	"github.com/blog/pkg/app"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	server *http.Server
}

func (s *HttpServer) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.server.Shutdown(ctx)
	}()
	return s.server.ListenAndServe()
}

func NewHttpServer(service *service.BlogService) app.HttpServer {
	server := new(HttpServer)
	engine := gin.Default()
	//加载api proto HTTPServer
	//v1.RegisterBlogHTTPServer(engine, service)
	server.server = &http.Server{
		Addr:    ":8888",
		Handler: engine,
	}
	return server
}
