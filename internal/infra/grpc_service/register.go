package grpc_service

import (
	"go.uber.org/dig"
	"google.golang.org/grpc"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment/v1"
	"github.com/suzushin54/experimental-parallel-api/internal/service"
)

func RegisterServices(s *grpc.Server, container *dig.Container) error {
	//return container.Invoke(func(paymentService *service.ParallelPaymentService) {
	//	pb.RegisterPaymentServiceServer(s, paymentService)
	//})
	return container.Invoke(func(paymentService *service.EventSourcedPaymentService) {
		pb.RegisterPaymentServiceServer(s, paymentService)
	})
}
