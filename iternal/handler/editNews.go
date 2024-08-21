package handler

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/LaughG33k/TestTask2/iternal"
	"github.com/LaughG33k/TestTask2/iternal/model"
	"github.com/LaughG33k/TestTask2/iternal/repository"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type EditNews struct {
	FiberApp      *fiber.App
	Repo          *repository.News
	JwtWorker     *iternal.JwtWorker
	Ctx           context.Context
	RquestTiemout time.Duration
}

func (h *EditNews) Handle() {

	h.FiberApp.Post("/edit/:id", h.hanleQuery, func(c fiber.Ctx) error {

		val, ok := c.GetReqHeaders()["Authorization"]

		if !ok {
			c.Status(403)
			return errors.New("not authorized")
		}

		token := strings.Split(val[0], " ")[1]

		m, err := h.JwtWorker.ParseJwt(token)

		if err != nil {
			c.Status(403)
			return errors.New("not authorized")
		}

		if m["permission"] != "admin" {
			c.Status(403)
			return errors.New("not authorized")
		}

		c.Next()
		return nil
	})

}

func (h *EditNews) hanleQuery(c fiber.Ctx) error {

	newsModel := model.NewsModel{}
	c.Accepts("application/json")

	idInt, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		c.Status(400)
		return errors.New("bad request")
	}

	if err := json.Unmarshal(c.Body(), &newsModel); err != nil {
		c.Status(500)
		return errors.New("iternal error. Decode error")
	}

	tm, canc := context.WithTimeout(h.Ctx, h.RquestTiemout)

	defer canc()

	if err := h.Repo.Edit(tm, int64(idInt), newsModel); err != nil {
		c.Status(500)
		return err
	}

	c.Status(200)

	return nil
}
