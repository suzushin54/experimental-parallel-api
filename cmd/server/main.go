package main

import (
	"log"
	"net"

	"github.com/suzushin54/experimental-parallel-api/internal/config"
	"github.com/suzushin54/experimental-parallel-api/internal/server"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal("failed to load configuration:", err)
	}

	lis, err := net.Listen("tcp", cfg.Server.Address)
	if err != nil {
		log.Fatal("failed to listen on port:", err)
	}

	s := grpc.NewServer()
	server.RegisterServices(s)

	log.Printf("server listening on %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
