package main

import (
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	appServer := wireApp()
	appServer.Run(ctx)
	defer cancel()
}
