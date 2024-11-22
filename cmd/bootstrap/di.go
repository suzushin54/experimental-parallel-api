package bootstrap

import (
	"log"

	"github.com/suzushin54/experimental-parallel-api/internal/domain/port"
	repositoryinterface "github.com/suzushin54/experimental-parallel-api/internal/domain/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/adapter"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/gateway"
	"github.com/suzushin54/experimental-parallel-api/internal/infra/repository"
	"github.com/suzushin54/experimental-parallel-api/internal/service"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	if err := c.Provide(func() repositoryinterface.PaymentRepository {
		return repository.NewMemoryPaymentRepository()
	}); err != nil {
		log.Fatalf("failed to provide PaymentRepository: %v", err)
	}
	if err := c.Provide(func() port.Mailer {
		return adapter.NewMailer()
	}); err != nil {
		log.Fatalf("failed to provide Mailer: %v", err)
	}
	if err := c.Provide(gateway.NewIDaaSGateway); err != nil {
		log.Fatalf("failed to provide IDaaS gateway: %v", err)
	}
	if err := c.Provide(gateway.NewPaymentGateway); err != nil {
		log.Fatalf("failed to provide PaymentGateway: %v", err)
	}
	if err := c.Provide(service.NewParallelPaymentService); err != nil {
		log.Fatalf("failed to provide ParallelPaymentService: %v", err)
	}

	return c
}
