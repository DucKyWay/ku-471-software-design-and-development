# Example Problems Day 1–3

## Day 1 พื้นฐานภาษา Go

### 1.8.1 Order Total

```go
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
```

### 1.8.2 Slice Aliasing

```go
package main

import "fmt"

func originalDiscountAll(items []float64) []float64 {
	for i := range items {
		items[i] *= 0.9
	}
	return items
}

func discountAll(items []float64) []float64 {
	result := make([]float64, len(items))
	for i, price := range items {
		result[i] = price * 0.9
	}
	return result
}

func main() {
	originalPrices := []float64{100, 200, 300, 400}
	originalWeekend := originalPrices[0:2]
	originalDiscounted := originalDiscountAll(originalWeekend)
	fmt.Println("original aliasing result:")
	fmt.Println(originalPrices)     // [90 180 300 400]
	fmt.Println(originalWeekend)    // [90 180]
	fmt.Println(originalDiscounted) // [90 180]

	prices := []float64{100, 200, 300, 400}
	weekend := prices[0:2]
	discounted := discountAll(weekend)

	fmt.Println("fixed result:")
	fmt.Println(prices)     // [100 200 300 400]
	fmt.Println(weekend)    // [100 200]
	fmt.Println(discounted) // [90 180]

	a := prices[0:2]
	a = append(a, 999)
	fmt.Println(prices) // [100 200 999 400]

	b := prices[0:2:2]
	b = append(b, 888)
	fmt.Println(prices) // [100 200 999 400]
	fmt.Println(b)      // [100 200 888]
}
```

### 1.8.3 Repository and Dependency Injection

```go
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
```

## Day 2 HTTP, Clean Architecture และฐานข้อมูล

### 2.6.1 Cancel Endpoint ที่เขียนผิดหลัก

```go
package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrNotCancelable = errors.New("order cannot be cancelled in its current status")
	ErrUnknownStatus = errors.New("unknown order status")
)
```

**ไฟล์:** `day02/food-api/domain/order.go`

```go
package domain

type Status string

const (
	Pending    Status = "pending"
	Confirmed  Status = "confirmed"
	Packing    Status = "packing"
	Cooking    Status = "cooking"
	Delivering Status = "delivering"
	Completed  Status = "completed"
	Cancelled  Status = "cancelled"
)

var flow = []Status{Pending, Confirmed, Packing, Cooking, Delivering, Completed}

type OrderLine struct {
	Item MenuItem `json:"item"`
	Qty  int      `json:"qty"`
}

func (line OrderLine) Total() float64 { return line.Item.Price * float64(line.Qty) }

type Order struct {
	ID     string      `json:"id"`
	Status Status      `json:"status"`
	Lines  []OrderLine `json:"lines"`
}

func (order Order) Total() float64 {
	total := 0.0
	for _, line := range order.Lines {
		total += line.Total()
	}
	return total
}

func (order Order) CanCancel() bool {
	return order.Status == Pending || order.Status == Confirmed
}

func (order Order) CanAdvanceTo(next Status) bool {
	for i, status := range flow {
		if status == order.Status {
			return i+1 < len(flow) && flow[i+1] == next
		}
	}
	return false
}

func (status Status) Valid() bool {
	switch status {
	case Pending, Confirmed, Packing, Cooking, Delivering, Completed, Cancelled:
		return true
	default:
		return false
	}
}
```

**ไฟล์:** `day02/food-api/usecase/order.go`

```go
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
```

**ไฟล์:** `day02/food-api/delivery/http/order.go`

```go
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
```

### 2.6.2 Endpoint GET /api/orders?status=pending

```go
package domain

type MenuItem struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
```

**ไฟล์:** `day02/food-api/repository/memory/order.go`

```go
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
```

```go
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"cs471/example-problem-solutions/day02/food-api/domain"
)

type OrderRepo struct{ DB *sql.DB }

func (r OrderRepo) FindOrder(id string) (domain.Order, error) {
	var order domain.Order
	err := r.DB.QueryRowContext(context.Background(),
		`SELECT id, status FROM orders WHERE id = $1`, id,
	).Scan(&order.ID, &order.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	rows, err := r.DB.QueryContext(context.Background(), `
		SELECT m.id, m.name, m.price, l.qty
		FROM order_lines l
		JOIN menu_items m ON m.id = l.menu_id
		WHERE l.order_id = $1 ORDER BY m.id`, id)
	if err != nil {
		return domain.Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var line domain.OrderLine
		if err := rows.Scan(&line.Item.ID, &line.Item.Name, &line.Item.Price, &line.Qty); err != nil {
			return domain.Order{}, err
		}
		order.Lines = append(order.Lines, line)
	}
	return order, rows.Err()
}

func (r OrderRepo) Save(order domain.Order) error {
	result, err := r.DB.ExecContext(context.Background(),
		`UPDATE orders SET status = $2 WHERE id = $1`, order.ID, string(order.Status))
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r OrderRepo) ListOrders(status domain.Status) ([]domain.Order, error) {
	rows, err := r.DB.QueryContext(context.Background(),
		`SELECT id FROM orders WHERE status = $1 ORDER BY id`, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	orders := make([]domain.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindOrder(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
```

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	deliveryhttp "cs471/example-problem-solutions/day02/food-api/delivery/http"
	"cs471/example-problem-solutions/day02/food-api/domain"
	"cs471/example-problem-solutions/day02/food-api/repository/memory"
	"cs471/example-problem-solutions/day02/food-api/usecase"
	"github.com/gofiber/fiber/v2"
)

func main() {
	serve := flag.Bool("serve", false, "run the HTTP server on :3002")
	flag.Parse()

	repo := memory.NewOrderRepo([]domain.Order{
		{ID: "o1", Status: domain.Pending, Lines: []domain.OrderLine{{Item: domain.MenuItem{ID: "m1", Name: "Pad Thai", Price: 80}, Qty: 2}}},
		{ID: "o2", Status: domain.Confirmed, Lines: []domain.OrderLine{{Item: domain.MenuItem{ID: "m2", Name: "Tom Yum", Price: 120}, Qty: 1}}},
		{ID: "o3", Status: domain.Cooking, Lines: []domain.OrderLine{{Item: domain.MenuItem{ID: "m3", Name: "Green Curry", Price: 95}, Qty: 2}}},
		{ID: "o4", Status: domain.Packing, Lines: []domain.OrderLine{{Item: domain.MenuItem{ID: "m4", Name: "Thai Tea", Price: 45}, Qty: 1}}},
	})
	handler := deliveryhttp.NewHandler(usecase.NewOrderUseCase(repo))

	app := fiber.New()
	api := app.Group("/api")
	api.Get("/orders", handler.ListOrders)
	api.Post("/orders/:id/cancel", handler.CancelOrder)

	if *serve {
		log.Println("Day 2 example server listening on http://localhost:3002")
		log.Fatal(app.Listen(":3002"))
	}

	showResponse(app, http.MethodGet, "/api/orders?status=pending")
	showResponse(app, http.MethodPost, "/api/orders/o1/cancel")
	showResponse(app, http.MethodGet, "/api/orders?status=cancelled")
	showResponse(app, http.MethodPost, "/api/orders/o3/cancel")
	showResponse(app, http.MethodPost, "/api/orders/o4/cancel")
	showResponse(app, http.MethodGet, "/api/orders?status=banana")
}

func showResponse(app *fiber.App, method, path string) {
	request := httptest.NewRequest(method, path, nil)
	response, err := app.Test(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s -> %d %s\n", method, path, response.StatusCode, string(body))
}
```

### 2.6.3 Schema That Can Change

```sql
ALTER TABLE orders ADD COLUMN receiver_name TEXT;
ALTER TABLE orders ADD COLUMN receiver_phone TEXT;

CREATE INDEX idx_orders_receiver_phone ON orders(receiver_phone);
```

## Day 3 Authorization และ Rate Limit

### 3.13.1 Endpoint ที่ตรวจสิทธิ์ผิดชั้น

```go
package main

import "fmt"

type Role string

const (
	Customer Role = "customer"
	Staff    Role = "staff"
	Manager  Role = "manager"
)

type User struct {
	ID   string
	Role Role
}

func (u User) IsStaff() bool      { return u.Role == Staff }
func (u User) CanCancelAny() bool { return u.Role == Staff || u.Role == Manager }

type Status string

const (
	Pending   Status = "pending"
	Confirmed Status = "confirmed"
	Cooking   Status = "cooking"
)

type Order struct {
	UserID string
	Status Status
}

func (o Order) OwnedBy(id string) bool { return o.UserID == id }
func (o Order) CanCancel() bool        { return o.Status == Pending || o.Status == Confirmed }

func (o Order) CanBeCancelledBy(actor User) bool {
	if !o.CanCancel() {
		return false
	}
	return actor.CanCancelAny() || o.OwnedBy(actor.ID)
}

func (o Order) CanChangeStatusBy(actor User) bool { return actor.IsStaff() }

func main() {
	alice := User{ID: "u1", Role: Customer}
	bob := User{ID: "u2", Role: Customer}
	manager := User{ID: "u3", Role: Manager}
	order := Order{UserID: alice.ID, Status: Pending}

	fmt.Printf("owner can cancel: %t\n", order.CanBeCancelledBy(alice))
	fmt.Printf("other customer can cancel: %t\n", order.CanBeCancelledBy(bob))
	fmt.Printf("manager can cancel another user's order: %t\n", order.CanBeCancelledBy(manager))
	fmt.Printf("manager can change status: %t\n", order.CanChangeStatusBy(manager))
	fmt.Printf("manager can cancel after cooking starts: %t\n", (Order{UserID: alice.ID, Status: Cooking}).CanBeCancelledBy(manager))
}
```

```sql
ALTER TABLE users DROP CONSTRAINT users_role_check;
ALTER TABLE users
    ADD CONSTRAINT users_role_check
    CHECK (role IN ('customer', 'staff', 'manager'));
```

### 3.13.2 ออเดอร์ของฉัน

```go
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
```

### 3.13.3 Rate Limit ที่แยกตามผู้ใช้

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type actorKey struct{}

type User struct{ ID string }

func withActor(r *http.Request, user User) *http.Request {
	return r.WithContext(contextWithActor(r, user))
}

func contextWithActor(r *http.Request, user User) context.Context {
	return context.WithValue(r.Context(), actorKey{}, user)
}

type window struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*window
	limit  int
	period time.Duration
	keyOf  func(*http.Request) string
}

func NewRateLimiter(limit int, period time.Duration, keyOf func(*http.Request) string) *RateLimiter {
	if keyOf == nil {
		keyOf = func(r *http.Request) string { return r.RemoteAddr }
	}
	return &RateLimiter{
		hits: make(map[string]*window), limit: limit, period: period, keyOf: keyOf,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.keyOf(r)
		now := time.Now()
		rl.mu.Lock()
		current, ok := rl.hits[key]
		if !ok || !now.Before(current.resetAt) {
			current = &window{resetAt: now.Add(rl.period)}
			rl.hits[key] = current
		}
		current.count++
		count := current.count
		rl.mu.Unlock()

		if count > rl.limit {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("rate limit exceeded\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func status(handler http.Handler, req *http.Request) int {
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res.Code
}

func main() {
	endpoint := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	perUser := NewRateLimiter(1, time.Minute, func(r *http.Request) string {
		actor, _ := r.Context().Value(actorKey{}).(User)
		return actor.ID
	})
	protected := perUser.Limit(endpoint)
	base := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	base.RemoteAddr = "10.0.0.5:1234"

	fmt.Println("authenticated requests from the same IP, separate user keys:")
	fmt.Println("alice:", status(protected, withActor(base, User{ID: "alice"})))
	fmt.Println("bob:  ", status(protected, withActor(base, User{ID: "bob"})))
	fmt.Println("alice:", status(protected, withActor(base, User{ID: "alice"})))

	loginByIP := NewRateLimiter(1, time.Minute, nil).Limit(endpoint)
	fmt.Println("login requests share the IP key:")
	fmt.Println("attempt 1:", status(loginByIP, base))
	fmt.Println("attempt 2:", status(loginByIP, base))
}
```