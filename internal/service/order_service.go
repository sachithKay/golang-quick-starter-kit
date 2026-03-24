package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sachithKay/ghost/internal/repository"
)

// OrderService contains the business logic for orders
type OrderService interface {
	ProcessNewOrder(ctx context.Context, customerID string, amount float64) (*repository.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) ProcessNewOrder(ctx context.Context, customerID string, amount float64) (*repository.Order, error) {
	// 1. Business Logic Rule: Prevent negative amounts
	if amount <= 0 {
		return nil, fmt.Errorf("invalid order amount: %f", amount)
	}

	// 2. Construct the Order entity
	// In a real app, use a UUID generator like google/uuid for ID
	newOrder := &repository.Order{
		ID:         fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		CustomerID: customerID,
		Amount:     amount,
		Status:     "PENDING",
	}

	// 3. Save to database using the repository
	err := s.repo.CreateOrder(ctx, newOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return newOrder, nil
}

func (s *orderService) runVerificationChecks(ctx context.Context, order *repository.Order) (bool, error) {
	// 1. Check if customer exists (Call Customer Service)
	// 2. Check if payment method is valid
	// 3. Check if inventory is available
	defer slog.Info("Order verification completed", "order_id", order.ID)

	validationChannel := make(chan bool, 3)

	go func(ch chan<- bool) {
		ch <- true
	}(validationChannel)

	go func(ch chan<- bool) {
		ch <- true
	}(validationChannel)

	go func(ch chan<- bool) {
		ch <- true
	}(validationChannel)

	for i := 0; i < 3; i++ {
		result := <-validationChannel
		if !result {
			return false, fmt.Errorf("order verification failed")
		}
	}
	return true, nil
}
