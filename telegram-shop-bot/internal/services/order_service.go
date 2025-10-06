package services

import (
	"database/sql"
	"telegram-shop-bot/internal/models"
)

type OrderService struct {
	db *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) Create(userID int64, productID, quantity int, total float64) (*models.Order, error) {
	result, err := s.db.Exec(`
		INSERT INTO orders (user_id, product_id, quantity, total, status)
		VALUES (?, ?, ?, ?, ?)
	`, userID, productID, quantity, total, models.StatusPending)

	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Order{
		ID:        int(id),
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Total:     total,
		Status:    string(models.StatusPending),
	}, nil
}

func (s *OrderService) GetUserOrders(userID int64) ([]models.Order, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, product_id, quantity, status, total, created_at
		FROM orders 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 20
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Status, &o.Total, &o.CreatedAt)
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *OrderService) UpdateStatus(orderID int, status models.OrderStatus) error {
	_, err := s.db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, orderID)
	return err
}
