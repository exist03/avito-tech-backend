package handlers

import (
	"avito-tech-backend/domain"
	"avito-tech-backend/internal"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
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

//	func (hl *HandlersLogger) GetL(c *fiber.Ctx) error {
//		hl.l.Info().Msg("sending request")
//		resp, err := hl.h.Get(c)
//		if err != nil {
//			hl.l.Warn().Msg("zalupa")
//			return err
//		}
//		hl.l.Info().Msg(string(resp))
//		c.Send(resp)
//		return nil
//	}

// TODO fix NoContent/PgxNoRows
// @Description Get list of user`s segments.
// @Summary Get list
// @Tags User
// @Produce json
// @Param user_id path int true "Segment ID"
// @Success 200 "OK"
// @Success 204 "No content"
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
// @Router /api/service/user/get/{user_id} [get]
func (h *Handlers) Get(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	response, err := h.service.Get(userId)
	if err != nil {
		if errors.Is(err, domain.ErrNoContent) {
			return c.SendStatus(http.StatusNoContent)
		} else if errors.Is(err, domain.ErrNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}
	}

	return c.Send(response)
}

// @Description Get file with history of user`s segments
// @Summary get file with history
// @Tags User
// @Accept json
// @Produce json
// @Param user_id query int true "Segment ID"
// @Param start query int true "Time start"
// @Param end query int true "Time end"
// @Success 200 "OK"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/service/user/get_history [get]
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

// @Description Create a new segment.
// @Summary create a new segment
// @Tags Segment
// @Accept json
// @Produce json
// @Param User-role header string false "admin"
// @Param segment_attrs body internal.Segment false "Segment attributes"
// @Success 200 "OK"
// @Failure 403 "Forbidden"
// @Failure 400 "Bad request"
// @Router /api/service/segment [post]
func (h *Handlers) Create(c *fiber.Ctx) error {
	segment := internal.Segment{}
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
		log.Println(err)
	}
	return c.SendStatus(http.StatusOK)
}

// @Description Update user`s segments.
// @Summary update user
// @Tags User
// @Accept json
// @Produce json
// @Param segment_attrs body internal.UpdateRequest false "Update attributes"
// @Success 200 "OK"
// @Failure 400 "Bad request"
// @Failure 403 "Forbidden"
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /api/service/user/update [patch]
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

// @Description Delete a segment.
// @Summary delete a new segment
// @Tags Segment
// @Accept json
// @Produce json
// @Param id path int true "Segment ID"
// @Param User-role header string false "admin"
// @Success 200 "OK"
// @Failure 404 "Not found"
// @Failure 400 "Bad request"
// @Router /api/service/segment/{id} [delete]
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
