//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/blog/internal/server"
	"github.com/blog/internal/service"
	"github.com/blog/pkg/app"
	"github.com/google/wire"
)

func wireApp() *app.App {
	panic(wire.Build(server.ProvideSet, service.ProvideSet, app.NewApp))
}
