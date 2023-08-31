package main

import (
	"avito-tech-backend/config"
	"avito-tech-backend/internal/app"
	"context"
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	cfg := config.GetConfigYml()
	a := app.New(ctx, cfg)
	app.Run(a, cfg)
}
