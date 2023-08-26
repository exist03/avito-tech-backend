package app

import (
	"avito-tech-backend/config"
	"avito-tech-backend/internal/handlers"
	"avito-tech-backend/internal/repository"
	"avito-tech-backend/internal/service"
	"context"
	"github.com/gofiber/fiber/v2"
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
	repository *repository.PsqlRepo // mb inject interface
	router     *fiber.App
}

func myMiddleware(handler func(c *fiber.Ctx) ([]byte, error), ctx *fiber.Ctx) fiber.Handler {
	log.Println("begin")
	data, _ := handler(ctx)
	return func(ctx *fiber.Ctx) error {
		ctx.Send(data)
		return nil
	}
	//log.Println("end")
}

func New(ctx context.Context, cfg *config.Config) *App {
	a := &App{}
	a.repository = repository.NewPsql(ctx, cfg.PsqlStorage)
	a.service = service.New(a.repository)
	a.handlers = handlers.New(a.service)
	a.router = fiber.New()

	//hl := handlers.NewHL(a.handlers, logger.GetLogger())
	//a.router.Use(logger.New())

	//a.router.Get("/", a.GetUserIDLogger)

	//a.router.Get("/api/service/get/:user_id", hl.GetL)
	a.router.Get("/api/service/get/:user_id", a.handlers.Get)
	a.router.Get("/api/service/get_history/", a.handlers.GetHistory)
	a.router.Post("/api/service/segment/", a.handlers.Create)
	a.router.Patch("/api/service/update/", a.handlers.Update)
	a.router.Delete("/api/service/segment/:id", a.handlers.Delete)
	return a
}

func Run(a *App) {
	//auto remove segments
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
