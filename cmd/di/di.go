package di

import (
	"log"

	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
	"github.com/suzushin54/experimental-parallel-api/internal/service"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	//if err := c.Provide(repository.NewPaymentRepository); err != nil {
	//	log.Fatalf("failed to provide PaymentRepository: %v", err)
	//}
	if err := c.Provide(gateway.NewPaymentGateway); err != nil {
		log.Fatalf("failed to provide PaymentGateway: %v", err)
	}
	if err := c.Provide(service.NewPaymentService); err != nil {
		log.Fatalf("failed to provide PaymentService: %v", err)
	}

	return c
}
