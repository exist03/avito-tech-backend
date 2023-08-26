package handlers

import (
	"avito-tech-backend/domain"
	"avito-tech-backend/internal"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type Service interface {
	Get(userId int) ([]byte, error)
	GetHistory(timeBegin, timeEnd int64, userId int) (string, error)
	Create(segment internal.Segment) error
	Delete(segmentId int) error
	Update(req internal.UpdateRequest) error
}

type Handlers struct {
	service Service
}

//type HandlersLogger struct {
//	h *Handlers
//	l zerolog.Logger
//}
//
//func NewHL(handlers *Handlers, logger zerolog.Logger) *HandlersLogger {
//	return &HandlersLogger{
//		h: handlers,
//		l: logger,
//	}
//}

func New(service Service) *Handlers {
	return &Handlers{service: service}
}

//func (hl *HandlersLogger) GetL(c *fiber.Ctx) error {
//	hl.l.Info().Msg("sending request")
//	resp, err := hl.h.Get(c)
//	if err != nil {
//		hl.l.Warn().Msg("zalupa")
//		return err
//	}
//	hl.l.Info().Msg(string(resp))
//	c.Send(resp)
//	return nil
//}

// TODO add NotFound
func (h *Handlers) Get(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	response, err := h.service.Get(userId)
	if err != nil {
		if errors.Is(err, domain.ErrNoContent) {
			return c.SendStatus(http.StatusNoContent)
		}
	}

	return c.Send(response)
}

func (h *Handlers) GetHistory(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
	}
	timeStart, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
	}
	timeEnd, err := strconv.ParseInt(c.Query("end"), 10, 64)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
	}
	response, err := h.service.GetHistory(timeStart, timeEnd, userId)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidArgument) {
			return c.SendStatus(http.StatusBadRequest)
		} else if errors.Is(err, domain.ErrNoContent) {
			return c.SendStatus(http.StatusNoContent)
		} else if errors.Is(err, domain.ErrNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendFile(response)
}

func (h *Handlers) Create(c *fiber.Ctx) error {
	var segment internal.Segment
	err := c.BodyParser(&segment)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	err = h.service.Create(segment)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidArgument) {
			return c.SendStatus(http.StatusBadRequest)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handlers) Update(c *fiber.Ctx) error {
	var req internal.UpdateRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	err = h.service.Update(req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handlers) Delete(c *fiber.Ctx) error {
	segmentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	err = h.service.Delete(segmentId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}
