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
	repository *repository.PsqlRepo
	router     *fiber.App
}

func New(ctx context.Context, cfg *config.Config) *App {
	a := &App{}
	a.repository = repository.NewPsql(ctx, cfg.PsqlStorage)
	a.service = service.New(a.repository)
	a.handlers = handlers.New(a.service)
	a.router = fiber.New()
	/*
		Метод добавления пользователя в сегмент. Принимает список (названий) сегментов которые нужно
		добавить пользователю, список (названий) сегментов которые нужно удалить у пользователя, id пользователя.
	*/
	a.router.Get("/api/service/get/:user_id", a.handlers.Get)
	a.router.Post("/api/service/segment/", a.handlers.Create)
	a.router.Post("/api/service/update/", a.handlers.Update)
	a.router.Delete("/api/service/segment/:id", a.handlers.Delete)
	return a
}

func Run(a *App) {
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
