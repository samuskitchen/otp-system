package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	ClientID    string `json:"client_id"`
}

type OTPValidateRequest struct {
	ClientID    string `json:"client_id"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type OTPResponse struct {
	Code string `json:"code"`
}

type OTP struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ClientID    string             `bson:"client_id"`
	PhoneNumber string             `bson:"phone_number"`
	Code        string             `bson:"code"`
	CreatedAt   time.Time          `bson:"created_at"`
	ExpiresAt   time.Time          `bson:"expires_at"`
	Attempts    int                `bson:"attempts"`
	Obsolete    bool               `bson:"obsolete"`
}

type Blacklist struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ClientID    string             `bson:"client_id"`
	PhoneNumber string             `bson:"phone_number"`
	BlockedAt   string             `bson:"blocked_at"`
}

type FailedAttempts struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ClientID    string             `bson:"client_id"`
	Code        string             `bson:"code"`
	PhoneNumber string             `bson:"phone_number"`
	Attempts    int                `bson:"attempts"`
	FailedAt    time.Time          `bson:"failed_at"`
}

func (otp *OTP) OTPToResponse() OTPResponse {
	return OTPResponse{
		Code: otp.Code,
	}
}
