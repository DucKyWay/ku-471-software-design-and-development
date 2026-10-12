package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
)

type Role string

const (
	Customer Role = "customer"
	Staff    Role = "staff"
)

type User struct {
	ID   string
	Role Role
}

func (u User) IsStaff() bool { return u.Role == Staff }

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
	Total  int    `json:"total"`
}

type OrderRepository interface {
	ListOrdersByUser(userID string) ([]Order, error)
	ListAllOrders() ([]Order, error)
}

type MemoryRepo struct{ orders []Order }

func (r MemoryRepo) ListOrdersByUser(userID string) ([]Order, error) {
	orders := make([]Order, 0)
	for _, order := range r.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return sorted(orders), nil
}

func (r MemoryRepo) ListAllOrders() ([]Order, error) {
	return sorted(append([]Order(nil), r.orders...)), nil
}

func sorted(orders []Order) []Order {
	sort.Slice(orders, func(i, j int) bool { return orders[i].ID < orders[j].ID })
	return orders
}

type OrderUseCase struct{ repo OrderRepository }

func (u OrderUseCase) ListMine(actor User) ([]Order, error) {
	if actor.IsStaff() {
		return u.repo.ListAllOrders()
	}
	return u.repo.ListOrdersByUser(actor.ID)
}

type actorKey struct{}

func withActor(r *http.Request, actor User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), actorKey{}, actor))
}

type Handler struct{ uc OrderUseCase }

func (h Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	actor, ok := r.Context().Value(actorKey{}).(User)
	if !ok || actor.ID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	orders, err := h.uc.ListMine(actor)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(orders)
}

func call(handler http.Handler, actor User) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, withActor(req, actor))
	return res
}

func main() {
	repo := MemoryRepo{orders: []Order{
		{ID: "o1", UserID: "u1", Status: "pending", Total: 160},
		{ID: "o2", UserID: "u2", Status: "confirmed", Total: 120},
		{ID: "o3", UserID: "u1", Status: "cooking", Total: 220},
	}}
	handler := Handler{uc: OrderUseCase{repo: repo}}
	route := http.HandlerFunc(handler.ListOrders)

	for _, actor := range []User{
		{ID: "u1", Role: Customer},
		{ID: "u2", Role: Customer},
		{ID: "u9", Role: Staff},
	} {
		res := call(route, actor)
		fmt.Printf("%s (%s): %d %s", actor.ID, actor.Role, res.Code, res.Body.String())
	}
}
