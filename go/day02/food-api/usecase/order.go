package usecase

import (
	"fmt"

	"cs471/example-problem-solutions/day02/food-api/domain"
)

type OrderRepository interface {
	FindOrder(id string) (domain.Order, error)
	Save(order domain.Order) error
	ListOrders(status domain.Status) ([]domain.Order, error)
}

type OrderUseCase struct{ repo OrderRepository }

func NewOrderUseCase(repo OrderRepository) OrderUseCase {
	return OrderUseCase{repo: repo}
}

func (u OrderUseCase) Cancel(id string) error {
	order, err := u.repo.FindOrder(id)
	if err != nil {
		return fmt.Errorf("cancel order %s: %w", id, err)
	}
	if !order.CanCancel() {
		return domain.ErrNotCancelable
	}
	order.Status = domain.Cancelled
	return u.repo.Save(order)
}

func (u OrderUseCase) ListOrders(status domain.Status) ([]domain.Order, error) {
	if !status.Valid() {
		return nil, domain.ErrUnknownStatus
	}
	return u.repo.ListOrders(status)
}
