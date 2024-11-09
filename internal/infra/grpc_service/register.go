package grpc_service

import (
	"google.golang.org/grpc"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment"
	"github.com/suzushin54/experimental-parallel-api/internal/service"
)

func RegisterServices(s *grpc.Server) {
	pb.RegisterPaymentServiceServer(s, service.NewPaymentService())
}
