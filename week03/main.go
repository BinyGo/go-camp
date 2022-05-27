package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type App struct {
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
		}
	}()

	// 带超时errgroup
	// ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
	// g, ctx := errgroup.WithContext(ctx)
	g, ctx := errgroup.WithContext(context.Background())

	//启动一个serveApp
	g.Go(func() error {
		app := serveApp{}
		return serve(ctx, ":8888", app.Start)
	})

	//启动一个serveDebug
	g.Go(func() error {
		debug := &serveDebug{}
		return serve(ctx, ":8999", debug.Start)
	})

	//监听 signal 信号处理
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			fmt.Println("signal quit:", sig)
			return fmt.Errorf("get os signal: %v", sig)
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group error: ", err)
		fmt.Println("ctx error: ", ctx.Err().Error())
	}
	//time.Sleep(time.Second)
	fmt.Println("all done!")

}

func serve(ctx context.Context, addr string, handler http.HandlerFunc) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		//监听ctx.cancel,让外面可以控制退出
		<-ctx.Done()
		fmt.Printf("server%s Shutdown \n", addr)
		err := s.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("server%s Shutdown err:%s\n", addr, err)
		}
		fmt.Printf("server%s Shutdown success\n", addr)
	}()
	return s.ListenAndServe()
}

type serveApp struct{}

func (srv *serveApp) Start(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("serveApp panic:", err)
		}
	}()
	req := fmt.Sprintf("hello serveApp: %+v\n", r)
	fmt.Fprint(w, req)
}

type serveDebug struct{}

func (srv *serveDebug) Start(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("serveDebug panic:", err)
		}
	}()
	req := fmt.Sprintf("hello serveDebug: %+v\n", r)
	fmt.Fprint(w, req)
}
