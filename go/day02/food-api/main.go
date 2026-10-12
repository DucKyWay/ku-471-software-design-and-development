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
