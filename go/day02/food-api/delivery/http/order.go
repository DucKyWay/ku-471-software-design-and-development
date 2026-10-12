package http

import (
	"errors"

	"cs471/example-problem-solutions/day02/food-api/domain"
	"cs471/example-problem-solutions/day02/food-api/usecase"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{ uc usecase.OrderUseCase }

func NewHandler(uc usecase.OrderUseCase) Handler { return Handler{uc: uc} }

func (h Handler) CancelOrder(c *fiber.Ctx) error {
	if err := h.uc.Cancel(c.Params("id")); err != nil {
		return toHTTP(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h Handler) ListOrders(c *fiber.Ctx) error {
	status := domain.Status(c.Query("status", string(domain.Pending)))
	orders, err := h.uc.ListOrders(status)
	if err != nil {
		return toHTTP(c, err)
	}
	return c.JSON(orders)
}

func toHTTP(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, domain.ErrNotCancelable):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, domain.ErrUnknownStatus):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}
}
