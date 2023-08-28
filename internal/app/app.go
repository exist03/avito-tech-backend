package app

import (
	"avito-tech-backend/config"
	_ "avito-tech-backend/docs"
	"avito-tech-backend/internal/handlers"
	"avito-tech-backend/internal/middleware"
	"avito-tech-backend/internal/repository"
	"avito-tech-backend/internal/service"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	httpPort = ":8080"
)

type App struct {
	handlers   *handlers.Handlers
	service    *service.Service
	repository *repository.PsqlRepo
	router     *fiber.App
}

func New(ctx context.Context, cfg *config.Config) *App {
	a := &App{}
	a.repository = repository.NewPsql(ctx, cfg.PsqlStorage)
	a.service = service.New(a.repository)
	a.handlers = handlers.New(a.service)
	a.router = fiber.New()
	a.router.Use(logger.New())
	a.router.Get("/swagger/*", swagger.HandlerDefault)
	//hl := handlers.NewHL(a.handlers, logger.GetLogger())
	//a.router.Use(logger.New())

	//a.router.Get("/", a.GetUserIDLogger)

	//a.router.Get("/api/service/get/:user_id", hl.GetL)

	a.router.Get("/api/service/user/get/:user_id", a.handlers.Get)
	a.router.Get("/api/service/user/get_history/", a.handlers.GetHistory)
	a.router.Post("/api/service/segment/", middleware.RoleCheck(a.handlers.Create))
	a.router.Patch("/api/service/user/update/", a.handlers.Update)
	a.router.Delete("/api/service/segment/:id", middleware.RoleCheck(a.handlers.Delete))
	return a
}

func Run(a *App) {
	//auto disconnect segments
	go a.repository.Checker()
	//Graceful	Shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Gracefully shutdown")
		if err := a.router.ShutdownWithTimeout(30 * time.Second); err != nil {
			log.Fatalln("server shutdown error: ", err)
		}
	}()

	err := a.router.Listen(httpPort)
	if err != nil {
		log.Fatalln(err)
	}
}
