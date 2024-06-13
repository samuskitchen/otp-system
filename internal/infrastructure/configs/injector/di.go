package injector

import (
	"fmt"

	"otp-system/internal/adapters/incoming/handler"
	"otp-system/internal/adapters/outgoing/repository"
	"otp-system/internal/core/service"
	"otp-system/internal/infrastructure/configs/router"
	"otp-system/internal/infrastructure/configs/router/otp"
	"otp-system/internal/infrastructure/configs/server"
	"otp-system/internal/infrastructure/configs/storage"

	"go.uber.org/dig"
)

var Container *dig.Container

func BuildContainer() *dig.Container {
	Container = dig.New()

	// DB
	checkError(Container.Provide(storage.ConnInstance))

	// Router Server
	checkError(Container.Provide(server.NewServer))
	checkError(Container.Provide(router.NewRouter))

	// Routers
	checkError(Container.Provide(otp.NewOTPRoute))

	// Handlers
	checkError(Container.Provide(handler.NewOTPHandler))

	// Service
	checkError(Container.Provide(service.NewOTPService))

	// Repository
	checkError(Container.Provide(repository.NewOtpRepository))

	return Container
}

func checkError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error injecting %v", err))
	}
}
