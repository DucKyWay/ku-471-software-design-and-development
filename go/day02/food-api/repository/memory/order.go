package memory

import (
	"sort"
	"sync"

	"cs471/example-problem-solutions/day02/food-api/domain"
)

type OrderRepo struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewOrderRepo(seed []domain.Order) *OrderRepo {
	repo := &OrderRepo{orders: make(map[string]domain.Order, len(seed))}
	for _, order := range seed {
		repo.orders[order.ID] = order
	}
	return repo
}

func (r *OrderRepo) FindOrder(id string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return domain.Order{}, domain.ErrNotFound
	}
	return order, nil
}

func (r *OrderRepo) Save(order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[order.ID]; !ok {
		return domain.ErrNotFound
	}
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepo) ListOrders(status domain.Status) ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]domain.Order, 0)
	for _, order := range r.orders {
		if order.Status == status {
			orders = append(orders, order)
		}
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].ID < orders[j].ID })
	return orders, nil
}
