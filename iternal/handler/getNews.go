package handler

import (
	"context"
	"errors"
	"time"

	"github.com/LaughG33k/TestTask2/iternal/model"
	"github.com/LaughG33k/TestTask2/iternal/repository"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type GetNews struct {
	Fiber         *fiber.App
	Repo          *repository.News
	RequstTimeout time.Duration
	Ctx           context.Context
}

func (h *GetNews) Handle() {
	h.Fiber.Get("/list", h.handleQuery)
}

func (h *GetNews) handleQuery(c fiber.Ctx) error {

	c.Accepts("application/json")

	tm, canc := context.WithTimeout(h.Ctx, h.RequstTimeout)

	defer canc()

	news, err := h.Repo.GetNews(tm)

	if err != nil {
		c.Status(500)
		return errors.New("Iternal error. Cant get a news")
	}

	list := model.ListNews{Success: true, News: make([]model.NewsModel, len(news))}

	for i := 0; i < len(news); i++ {
		list.News[i] = *news[i]
	}

	bytes, err := json.Marshal(list)

	if err != nil {
		c.Status(500)
		return errors.New("Iternal")
	}

	c.Status(200)
	c.Send(bytes)

	return nil
}
