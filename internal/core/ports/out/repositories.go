package out

import (
	"context"

	"otp-system/internal/core/domain"
)

type OTPRepository interface {
	CheckBlacklist(ctx context.Context, phoneNumber string) (bool, error)
	InsertOTP(ctx context.Context, clientID, phoneNumber, codeOTP string) (domain.OTP, error)
	InsertBlackList(ctx context.Context, clientID, phoneNumber string) (domain.Blacklist, error)
	FindOTP(ctx context.Context, phoneNumber, code string) (*domain.OTP, error)
	UpdateOTP(ctx context.Context, otp *domain.OTP) error
	MarkOTPAsObsolete(ctx context.Context, phoneNumber, code string) error
	LogFailedAttempt(ctx context.Context, clientID, phoneNumber, code string) (int, error)
}
