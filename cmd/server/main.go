package main

import (
	"fmt"
	"log"
	"net"

	"github.com/suzushin54/experimental-parallel-api/cmd/config"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/grpc_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load configuration:", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal("failed to listen on port:", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	grpc_service.RegisterServices(s)

	log.Printf("grpc_service listening on %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
