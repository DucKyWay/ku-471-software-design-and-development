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
