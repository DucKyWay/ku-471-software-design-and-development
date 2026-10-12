package main

import (
	"errors"
	"fmt"
)

var (
	ErrBelowMinimum = errors.New("order below 100 baht")
	ErrInvalidQty   = errors.New("qty must be at least 1")
)

type MenuItem struct {
	ID    string
	Name  string
	Price float64
}

type OrderLine struct {
	Item MenuItem
	Qty  int
}

func (l OrderLine) Total() float64 {
	return l.Item.Price * float64(l.Qty)
}

func OrderTotal(lines []OrderLine, deliveryFee float64) (float64, error) {
	subtotal := 0.0
	for _, line := range lines {
		subtotal += line.Total()
	}

	if subtotal < 100 {
		return 0, ErrBelowMinimum
	}
	return subtotal + deliveryFee, nil
}

func NewOrderLine(item MenuItem, qty int) (OrderLine, error) {
	if qty < 1 {
		return OrderLine{}, ErrInvalidQty
	}
	return OrderLine{Item: item, Qty: qty}, nil
}

func BuildOrder(items []MenuItem, quantities []int) ([]OrderLine, error) {
	if len(items) != len(quantities) {
		return nil, fmt.Errorf("items and quantities have different lengths")
	}

	lines := make([]OrderLine, 0, len(items))
	for i, item := range items {
		line, err := NewOrderLine(item, quantities[i])
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func main() {
	lines := []OrderLine{
		{Item: MenuItem{Name: "Pad Thai", Price: 80}, Qty: 2},
		{Item: MenuItem{Name: "Tom Yum", Price: 120}, Qty: 1},
	}

	total, err := OrderTotal(lines, 30)
	if err != nil {
		if errors.Is(err, ErrBelowMinimum) {
			fmt.Println("ยอดสั่งซื้อขั้นต่ำ 100 บาท")
			return
		}
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%.2f\n", total)

	_, err = BuildOrder(
		[]MenuItem{{Name: "Pad Thai", Price: 80}, {Name: "Tom Yum", Price: 120}},
		[]int{2, 0},
	)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Is(err, ErrInvalidQty))
	}
}
