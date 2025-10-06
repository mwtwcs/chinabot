package services

import (
	"database/sql"
	"telegram-shop-bot/internal/models"
)

type ProductService struct {
	db *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, price, image, category, stock, created_at 
		FROM products 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Image, &p.Category, &p.Stock, &p.CreatedAt)
		if err != nil {
			continue
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRow(`
		SELECT id, name, description, price, image, category, stock, created_at 
		FROM products WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Image, &p.Category, &p.Stock, &p.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ProductService) Create(p *models.Product) error {
	result, err := s.db.Exec(`
		INSERT INTO products (name, description, price, image, category, stock)
		VALUES (?, ?, ?, ?, ?, ?)
	`, p.Name, p.Description, p.Price, p.Image, p.Category, p.Stock)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (s *ProductService) Update(p *models.Product) error {
	_, err := s.db.Exec(`
		UPDATE products 
		SET name = ?, description = ?, price = ?, image = ?, category = ?, stock = ?
		WHERE id = ?
	`, p.Name, p.Description, p.Price, p.Image, p.Category, p.Stock, p.ID)
	return err
}

func (s *ProductService) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}
