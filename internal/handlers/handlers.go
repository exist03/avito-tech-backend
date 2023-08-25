package handlers

import (
	"avito-tech-backend/internal"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
)

type Service interface {
	Get(userId int) ([]byte, error)
	Create(segment internal.Segment) error
	Delete(segmentId int) error
	Update(req internal.UpdateRequest) error
}

type Handlers struct {
	service Service
}

func New(service Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) Get(ctx *fiber.Ctx) error {
	userId, err := strconv.Atoi(ctx.Params("user_id"))
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
	}
	response, _ := h.service.Get(userId)

	return ctx.Send(response)
}

func (h *Handlers) Create(ctx *fiber.Ctx) error {
	var segment internal.Segment
	err := ctx.BodyParser(&segment)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	err = h.service.Create(segment)
	if err != nil {
		log.Println(err)
		//TODO
	}
	return ctx.SendStatus(http.StatusOK)
}

func (h *Handlers) Update(ctx *fiber.Ctx) error {
	var req internal.UpdateRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	h.service.Update(req)
	return nil
}

func (h *Handlers) Delete(ctx *fiber.Ctx) error {
	segmentId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	h.service.Delete(segmentId)
	return nil
}
