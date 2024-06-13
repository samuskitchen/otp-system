package otp

import (
	"otp-system/internal/adapters/incoming/handler"
	"otp-system/pkg/kit/enums"

	"github.com/labstack/echo/v4"
)

type otpRoute struct {
	otpHandler handler.OTPHandler
}

type Route interface {
	Resource(c *echo.Group)
}

func NewOTPRoute(otpHandler handler.OTPHandler) Route {
	return &otpRoute{
		otpHandler: otpHandler,
	}
}

func (or *otpRoute) Resource(c *echo.Group) {
	c.POST(enums.GenerateOTPPOST, or.otpHandler.RequestOTP)
	c.POST(enums.ValidateOTPPOST, or.otpHandler.VerifyOTP)
}
