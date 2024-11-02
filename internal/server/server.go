package server

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/suzushin54/experimental-parallel-api/gen/payment"
)

type paymentService struct {
	pb.UnimplementedPaymentServiceServer
}

func NewPaymentService() *paymentService {
	return &paymentService{}
}

func (s *paymentService) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("ProcessPayment: %v", req)

	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: "1234567890",
		Message:       "Payment processed successfully",
		ErrorMessage:  "",
	}, nil
}

func RegisterServices(s *grpc.Server) {
	pb.RegisterPaymentServiceServer(s, NewPaymentService())
}
