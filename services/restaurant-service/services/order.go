package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/repository"
)

type OrderService interface {
	PlaceOrder(order *models.Order, items []models.OrderItem) (int64, error)
	GetOrderStatus(orderID int64) (string, error)
	UpdateOrderStatus(orderID int64, status string) error
	GetOrder(orderID int64) (*models.Order, error)
}

type orderService struct {
	repo repository.OrderRepo
	db   *sql.DB
}

func NewOrderService(r repository.OrderRepo, db *sql.DB) OrderService {
	return &orderService{repo: r, db: db}
}

func (s *orderService) PlaceOrder(order *models.Order, items []models.OrderItem) (int64, error) {
	if order == nil || len(items) == 0 {
		return 0, errors.New("order and items required")
	}
	// create tx
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		// ensure rollback on panic
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	order.CreatedAt = timePtr(time.Now().UTC())
	order.UpdatedAt = timePtr(time.Now().UTC())

	orderID, err := s.repo.CreateOrderWithItems(tx, order, items)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	// Optionally compute order_number or other post-processing here, then update order row.
	// For example set order_number = "ORD"+orderID or use some generator.
	// We'll commit and return id.

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return orderID, nil
}

func (s *orderService) GetOrderStatus(orderID int64) (string, error) {
	return s.repo.GetOrderStatus(orderID)
}

func (s *orderService) UpdateOrderStatus(orderID int64, status string) error {
	// optional validation of status
	if status == "" {
		return errors.New("status required")
	}
	// More advanced: check valid transitions
	return s.repo.UpdateOrderStatus(orderID, status)
}

func (s *orderService) GetOrder(orderID int64) (*models.Order, error) {
	return s.repo.GetOrderByID(orderID)
}
