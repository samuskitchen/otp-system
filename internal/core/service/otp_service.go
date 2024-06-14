package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"otp-system/internal/core/domain"
	"otp-system/internal/core/ports/in"
	"otp-system/internal/core/ports/out"

	log "github.com/rs/zerolog/log"
)

const (
	MaxAttempts int = 3
)

type otpService struct {
	otpRepository out.OTPRepository
}

func NewOTPService(otpRepository out.OTPRepository) in.OTPService {
	return &otpService{
		otpRepository: otpRepository,
	}
}

func (os *otpService) GenerateOTP(ctx context.Context, data domain.OTPRequest) (domain.OTPResponse, error) {
	log.Info().Msg("implement Service GenerateOTP")

	codeOTP, err := generateCodeOTP()
	if err != nil {
		return domain.OTPResponse{}, err
	}

	dataOTP, err := os.otpRepository.InsertOTP(ctx, data.ClientID, data.PhoneNumber, codeOTP)
	if err != nil {
		return domain.OTPResponse{}, err
	}

	return dataOTP.OTPToResponse(), nil
}

func (os *otpService) ValidateOTP(ctx context.Context, data domain.OTPValidateRequest) (bool, error) {
	log.Info().Msg("implement Service ValidateOTP")

	otp, err := os.otpRepository.FindOTP(ctx, data.PhoneNumber, data.Code)
	if err != nil {
		// If OTP is not found, log the failed attempt
		attempts, logErr := os.otpRepository.LogFailedAttempt(ctx, data.ClientID, data.PhoneNumber, data.Code)
		if logErr != nil {
			return false, logErr
		}

		if attempts >= MaxAttempts {
			_, blacklistErr := os.otpRepository.InsertBlackList(ctx, data.ClientID, data.PhoneNumber)
			if blacklistErr != nil {
				return false, blacklistErr
			}
		}

		return false, fmt.Errorf("OTP not found or already used")
	}

	if otp.Obsolete {
		return false, fmt.Errorf("OTP obsolete")
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		otp.Obsolete = true
		updateErr := os.otpRepository.UpdateOTP(ctx, otp)
		if updateErr != nil {
			return false, updateErr
		}

		return false, fmt.Errorf("OTP expired")
	}

	otp.Attempts++
	updateErr := os.otpRepository.UpdateOTP(ctx, otp)
	if updateErr != nil {
		return false, updateErr
	}

	return !otp.Obsolete, nil
}

func (os *otpService) CheckBlacklist(ctx context.Context, dataValidate domain.OTPValidateRequest) (bool, error) {
	log.Info().Msg("implement Service CheckBlacklist")
	return os.otpRepository.CheckBlacklist(ctx, dataValidate.PhoneNumber)
}

func (os *otpService) InsertBlackListOTP(ctx context.Context, data domain.OTPRequest) error {
	log.Info().Msg("implement Service InsertBlackListOTP")

	_, err := os.otpRepository.InsertBlackList(ctx, data.ClientID, data.PhoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func generateCodeOTP() (string, error) {
	const otpChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var err error
	otp := make([]byte, 5)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			log.Error().Err(err).Msg("GenerateCodeOTP")
		}

		otp[i] = otpChars[num.Int64()]
	}

	return string(otp), err
}
