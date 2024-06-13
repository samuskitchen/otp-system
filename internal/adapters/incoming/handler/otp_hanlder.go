package handler

import (
	"net/http"

	"otp-system/internal/core/domain"
	"otp-system/internal/core/ports/in"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type otpHandler struct {
	otpService in.OTPService
}

type OTPHandler interface {
	RequestOTP(c echo.Context) error
	VerifyOTP(c echo.Context) error
}

func NewOTPHandler(otpService in.OTPService) OTPHandler {
	return &otpHandler{
		otpService: otpService,
	}
}

func (oh *otpHandler) RequestOTP(c echo.Context) error {
	ctx := c.Request().Context()

	var otpRequest domain.OTPRequest
	if err := c.Bind(&otpRequest); err != nil {
		return err
	}

	otp, err := oh.otpService.GenerateOTP(ctx, otpRequest)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error generating OTP"})
	}

	return c.JSON(http.StatusOK, otp)
}

func (oh *otpHandler) VerifyOTP(c echo.Context) error {
	ctx := c.Request().Context()

	var otpValidateRequest domain.OTPValidateRequest
	if err := c.Bind(&otpValidateRequest); err != nil {
		return err
	}

	blacklisted, err := oh.otpService.CheckBlacklist(ctx, otpValidateRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if blacklisted {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "Phone number is blacklisted"})
	}

	valid, err := oh.otpService.ValidateOTP(ctx, otpValidateRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if !valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP verified successfully"})
}
