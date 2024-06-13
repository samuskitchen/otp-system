package router

import (
	"os"
	"strings"

	"otp-system/internal/infrastructure/configs/router/otp"
	customMiddleware "otp-system/internal/infrastructure/middleware"
	"otp-system/pkg/kit/enums"

	"github.com/labstack/echo/v4"
	middlewareEcho "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Router struct {
	server *echo.Echo
	otp    otp.Route
}

func NewRouter(server *echo.Echo, otp otp.Route) *Router {
	return &Router{
		server,
		otp,
	}
}

func (r *Router) Init() {
	// Custom zerolog logger instance
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Middleware
	logConfig := customMiddleware.ZeroLogConfig{
		Logger: logger,
		FieldMap: map[string]string{
			"uri":    "@uri",
			"host":   "@host",
			"method": "@method",
			"status": "@status",
		},
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, enums.HealthPath)
		},
	}

	r.server.Use(customMiddleware.ZeroLogWithConfig(logConfig))
	r.server.Use(middlewareEcho.Recover())
	r.server.Use(middlewareEcho.RequestID())

	apiGroup := r.server.Group(enums.BasePath)
	r.otp.Resource(apiGroup)

	//apiGroup.GET("/docs/*", swagger.EchoWrapHandler())
	//apiGroup.GET(enums.HealthPath, r.healthHandler.HealthCheck)

	for _, route := range r.server.Routes() {
		log.Info().Msgf("[%s] %s", route.Method, route.Path)
	}
}
