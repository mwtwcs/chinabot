package services

import (
	"database/sql"
	"telegram-shop-bot/internal/models"
)

type CartService struct {
	db *sql.DB
}

func NewCartService(db *sql.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) AddItem(userID int64, productID, quantity int) error {
	_, err := s.db.Exec(`
		INSERT INTO cart (user_id, product_id, quantity) 
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, product_id) 
		DO UPDATE SET quantity = quantity + ?
	`, userID, productID, quantity, quantity)
	return err
}

func (s *CartService) GetItems(userID int64) ([]models.CartItem, error) {
	rows, err := s.db.Query(`
		SELECT user_id, product_id, quantity 
		FROM cart 
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		rows.Scan(&item.UserID, &item.ProductID, &item.Quantity)
		items = append(items, item)
	}
	return items, nil
}

func (s *CartService) Clear(userID int64) error {
	_, err := s.db.Exec("DELETE FROM cart WHERE user_id = ?", userID)
	return err
}

func (s *CartService) RemoveItem(userID int64, productID int) error {
	_, err := s.db.Exec("DELETE FROM cart WHERE user_id = ? AND product_id = ?", userID, productID)
	return err
}
