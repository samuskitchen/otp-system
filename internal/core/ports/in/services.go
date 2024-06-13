package in

import (
	"context"

	"otp-system/internal/core/domain"
)

type OTPService interface {
	GenerateOTP(ctx context.Context, data domain.OTPRequest) (domain.OTPResponse, error)
	ValidateOTP(ctx context.Context, dataValidate domain.OTPValidateRequest) (bool, error)
	CheckBlacklist(ctx context.Context, dataValidate domain.OTPValidateRequest) (bool, error)
	InsertBlackListOTP(ctx context.Context, data domain.OTPRequest) error
}
