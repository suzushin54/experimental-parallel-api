package grpc_service

import (
	"go.uber.org/dig"
	"google.golang.org/grpc"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment"
	"github.com/suzushin54/experimental-parallel-api/internal/service"
)

func RegisterServices(s *grpc.Server, container *dig.Container) error {
	return container.Invoke(func(paymentService *service.PaymentService) {
		pb.RegisterPaymentServiceServer(s, paymentService)
	})
}
