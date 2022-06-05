package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server interface {
	Serve(ctx context.Context) error
}

type GrpcServer Server
type HttpServer Server

type App struct {
	http HttpServer
	grpc GrpcServer
}

func NewApp(http HttpServer, grpc GrpcServer) *App {
	return &App{
		http: http,
		grpc: grpc,
	}
}

func (a *App) Run(ctx context.Context) {
	group := new(errgroup.Group)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//启动httpServer
	group.Go(func() error {
		err := a.http.Serve(ctx)
		if err != nil {
			cancel()
		}
		return err
	})
	//启动grpcServer
	group.Go(func() error {
		err := a.grpc.Serve(ctx)
		if err != nil {
			cancel()
		}
		return err
	})
	//监听 signal 信号处理
	group.Go(func() error {
		err := ReceiveSignal(ctx)
		if err != nil {
			cancel()
		}
		return err
	})

	if err := group.Wait(); err != nil {
		fmt.Println("group error: ", err)
		fmt.Println("ctx error: ", ctx.Err().Error())
	}
}

func ReceiveSignal(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		fmt.Println("signal quit:", sig)
		return errors.New(sig.String())
	case <-ctx.Done():
		return ctx.Err()
	}
}
