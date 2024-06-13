package main

import (
	"fmt"
	"os"

	"otp-system/internal/infrastructure/configs/injector"
	"otp-system/internal/infrastructure/configs/router"
	"otp-system/internal/infrastructure/middleware/log"
	"otp-system/pkg/kit/enums"

	"github.com/labstack/echo/v4"
)

func main() {
	container := injector.BuildContainer()

	log.InitLogger(enums.App)

	err := container.Invoke(func(server *echo.Echo, route *router.Router) {
		address := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
		server.Debug = os.Getenv("SERVER_POSTFIX") == enums.PostfixDev

		route.Init()
		server.Logger.Fatal(server.Start(address))
	})

	if err != nil {
		panic(err)
	}
}
