package handler

import (
	"errors"

	"github.com/deblinux123/Redis-DB/day-07/internal/model"
	"github.com/deblinux123/Redis-DB/day-07/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invlaid JSON body",
		})
	}

	if err := h.service.Create(c.UserContext(), user); err != nil {
		if errors.Is(err, service.ErrInvlaidUser) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Ok",
		"User":   user,
	})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.service.GetUserByID(c.UserContext(), id)

	if err != nil {
		if errors.Is(err, service.ErrInvlaidUser) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if errors.Is(err, redis.Nil) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No User Found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get user.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Ok",
		"User":   user,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON body",
		})
	}

	user.ID = c.Params("id")

	if err := h.service.Update(c.UserContext(), user); err != nil {
		if errors.Is(err, service.ErrInvlaidUser) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Update Ok",
		"user":   user,
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, service.ErrInvlaidUser) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Deleted Ok",
	})
}
