package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("menu item not found")
	ErrBelowMinimum = errors.New("order below 100 baht")
)

type MenuItem struct {
	ID    string
	Name  string
	Price float64
}

type Repository interface {
	FindMenuItem(id string) (MenuItem, error)
}

type MemoryRepo struct {
	items map[string]MenuItem
}

func (r MemoryRepo) FindMenuItem(id string) (MenuItem, error) {
	item, ok := r.items[id]
	if !ok {
		return MenuItem{}, fmt.Errorf("find %s: %w", id, ErrNotFound)
	}
	return item, nil
}

type FakeRepo struct{}

func (FakeRepo) FindMenuItem(id string) (MenuItem, error) {
	return MenuItem{}, fmt.Errorf("find %s: %w", id, ErrNotFound)
}

type OrderUseCase struct {
	repo Repository
}

func (u OrderUseCase) Place(id string, qty int) (float64, error) {
	item, err := u.repo.FindMenuItem(id)
	if err != nil {
		return 0, err
	}

	total := item.Price * float64(qty)
	if total < 100 {
		return 0, ErrBelowMinimum
	}
	return total, nil
}

func main() {
	seed := map[string]MenuItem{
		"m1": {ID: "m1", Name: "Pad Thai", Price: 80},
	}

	useCase := OrderUseCase{repo: MemoryRepo{items: seed}}
	fmt.Println(useCase.Place("m1", 2)) // 160 <nil>

	failing := OrderUseCase{repo: FakeRepo{}}
	fmt.Println(failing.Place("m1", 2)) // 0 find m1: menu item not found
}
