package main

import (
	"avito-tech-backend/config"
	"avito-tech-backend/internal/app"
	"context"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfigYml()
	a := app.New(ctx, cfg)
	app.Run(a)
}
