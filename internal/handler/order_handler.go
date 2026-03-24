package handler

import (
	"context"

	"github.com/sachithKay/ghost/internal/config"
	"github.com/sachithKay/ghost/internal/service"

	// 1. Import your generated code
	orderv1 "github.com/sachithKay/ghost/gen/go/v1"
)

// orderHandler is the private implementation of the gRPC server interface
type orderHandler struct {
	// 2. EMBED the generated server to satisfy the gRPC interface
	orderv1.UnimplementedOrderServiceServer
	Config  *config.Config
	Service service.OrderService
}

// NewOrderHandler returns the public gRPC interface, hiding the implementation
func NewOrderHandler(cfg *config.Config, svc service.OrderService) orderv1.OrderServiceServer {
	return &orderHandler{Config: cfg, Service: svc}
}

// 3. Implement the RPC method defined in your .proto
// Note: We no longer need the Routes() method because the Gateway handles it!
func (h *orderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	// 1. Handlers call the Service (Business Logic)
	// We pass clean parameters, not GRPC structs, down to the service layer.
	order, err := h.Service.ProcessNewOrder(ctx, req.CustomerId, req.Amount)
	if err != nil {
		return nil, err
	}

	// 2. Handlers format the Response to match the protocol
	return &orderv1.CreateOrderResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}
