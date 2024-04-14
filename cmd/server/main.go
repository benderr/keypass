package main

import (
	"context"

	"github.com/benderr/keypass/internal/server/app"
	"github.com/benderr/keypass/internal/server/config"
)

func main() {
	conf := config.MustLoad()
	ctx := context.Background()

	app.Run(ctx, conf)
}
