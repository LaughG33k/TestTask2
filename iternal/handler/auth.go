package handler

import (
	"context"
	"errors"
	"time"

	"github.com/LaughG33k/TestTask2/iternal"
	"github.com/LaughG33k/TestTask2/iternal/model"
	"github.com/LaughG33k/TestTask2/iternal/repository"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type Auth struct {
	Fiber         *fiber.App
	JwtWorker     *iternal.JwtWorker
	Ctx           context.Context
	Repo          *repository.User
	RtRepo        *repository.RefreshToken
	RequstTimeout time.Duration
}

func (h *Auth) Handle() {

	h.Fiber.Post("/auth/register", h.register)
	h.Fiber.Post("/auth/login", h.loginIn)
	h.Fiber.Put("/auth/updatejwt", h.UpdateJwt)

}

func (h *Auth) register(c fiber.Ctx) error {

	c.Accepts("application/json")

	regModel := model.LoginPass{}

	if err := json.Unmarshal(c.Body(), &regModel); err != nil {
		c.Status(500)
		return errors.New("Iternal error. decode failed")
	}

	tm, canc := context.WithTimeout(h.Ctx, h.RequstTimeout)

	defer canc()

	if err := h.Repo.Register(tm, regModel.Login, regModel.Pass, "user"); err != nil {
		c.Status(500)
		return err
	}

	c.Status(200)

	return nil

}

func (h *Auth) loginIn(c fiber.Ctx) error {

	c.Accepts("application/json")

	loginModel := model.LoginPass{}

	if err := json.Unmarshal(c.Body(), &loginModel); err != nil {
		c.Status(500)
		return errors.New("Iternal error. decode failed")
	}

	tm, canc := context.WithTimeout(h.Ctx, h.RequstTimeout)

	defer canc()

	t, uuid, perm, err := h.Repo.CheckByLP(tm, loginModel.Login, loginModel.Pass)

	if err != nil {
		c.Status(500)

		return err
	}

	if !t {
		c.Status(404)

		return errors.New("User not found")
	}

	rt := iternal.GenerateRefreshToken(30)

	if err := h.RtRepo.CreateRefreshToken(tm, "", rt, uuid, time.Now().Add(15*24*time.Hour).Unix()); err != nil {
		c.Status(500)
		return errors.New("Iternal error")
	}

	jwtToken, err := h.JwtWorker.CreateJwt(map[string]any{"uuid": uuid, "permission": perm})

	if err != nil {
		c.Status(500)
		return errors.New("Iternal error")
	}

	bytes, err := json.Marshal(model.Tokens{Jwt: jwtToken, Refresh: rt})

	if err != nil {
		c.Status(500)
		return errors.New("Iternal error")
	}

	c.Status(200)
	c.Send(bytes)

	return nil

}

func (h *Auth) UpdateJwt(c fiber.Ctx) error {

	rt := map[string]string{"refresh_token": ""}

	c.Accepts("application/json")

	if err := json.Unmarshal(c.Body(), &rt); err != nil {
		c.Status(500)
		return errors.New("Iterabl error. decode failed")
	}

	if rt["refresh_token"] == "" {
		c.Status(400)
		return errors.New("Bad request")
	}

	tm, canc := context.WithTimeout(h.Ctx, h.RequstTimeout)

	defer canc()

	uuid, tokentimelive, err := h.RtRepo.FindRefreshToken(tm, rt["refresh_token"])

	if err != nil {
		c.Status(500)
		return err
	}

	if uuid == "" {
		c.Status(404)
		return errors.New("Not foud")
	}

	if time.Now().Unix() > tokentimelive {
		c.Status(401)
		return errors.New("token has expired")
	}

	newRt := iternal.GenerateRefreshToken(30)

	if err := h.RtRepo.CreateRefreshToken(tm, rt["refresh_token"], newRt, uuid, time.Now().Add(24*15*time.Hour).Unix()); err != nil {
		c.Status(500)
		return err
	}

	_, _, _, perm, err := h.Repo.GetAllFields(tm, uuid)

	if err != nil {
		c.Status(500)
		return err
	}

	jwtToken, err := h.JwtWorker.CreateJwt(map[string]any{"uuid": uuid, "permission": perm})

	if err != nil {
		c.Status(500)
		return err
	}

	bytes, err := json.Marshal(model.Tokens{Jwt: jwtToken, Refresh: newRt})

	if err != nil {
		c.Status(500)
		return err
	}

	c.Status(200)
	c.Send(bytes)

	return nil

}
