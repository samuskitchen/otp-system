package repository

import (
	"context"
	"time"

	"otp-system/internal/core/domain"
	"otp-system/internal/core/ports/out"
	"otp-system/pkg/kit/enums"

	log "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	OTPValidityDuration time.Duration = 1 * time.Minute
	MaxAttempts         int           = 3
)

type otpRepository struct {
	mongo *mongo.Client
}

func NewOtpRepository(mongoClient *mongo.Client) out.OTPRepository {
	return &otpRepository{
		mongo: mongoClient,
	}
}

func (or *otpRepository) InsertOTP(ctx context.Context, clientID, phoneNumber, codeOTP string) (domain.OTP, error) {
	log.Info().Msg("implement Insert")

	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginLogOTP)

	now := time.Now()
	createOtp := domain.OTP{
		PhoneNumber: phoneNumber,
		Code:        codeOTP,
		ClientID:    clientID,
		CreatedAt:   now,
		ExpiresAt:   now.Add(OTPValidityDuration),
		Attempts:    0,
	}

	result, err := collection.InsertOne(ctx, createOtp)
	if err != nil {
		return domain.OTP{}, err
	}

	createOtp.ID = result.InsertedID.(primitive.ObjectID)

	return createOtp, nil
}

func (or *otpRepository) InsertBlackList(ctx context.Context, clientID, phoneNumber string) (domain.Blacklist, error) {
	log.Info().Msg("implement Insert")

	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginBlackListOTP)

	dueDate := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	createBlacklist := domain.Blacklist{
		ClientID:    clientID,
		PhoneNumber: phoneNumber,
		BlockedAt:   dueDate,
	}

	result, err := collection.InsertOne(ctx, createBlacklist)
	if err != nil {
		return domain.Blacklist{}, err
	}

	createBlacklist.ID = result.InsertedID.(primitive.ObjectID)
	return createBlacklist, nil
}

func (or *otpRepository) FindOTP(ctx context.Context, phoneNumber, code string) (*domain.OTP, error) {
	log.Info().Msg("implement FindOTP")
	var otp domain.OTP

	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginLogOTP)
	err := collection.FindOne(ctx, bson.M{"phone_number": phoneNumber, "code": code}).Decode(&otp)

	return &otp, err
}

func (or *otpRepository) UpdateOTP(ctx context.Context, otp *domain.OTP) error {
	log.Info().Msg("implement UpdateOTP")

	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginLogOTP)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": otp})

	return err
}

func (or *otpRepository) CheckBlacklist(ctx context.Context, phoneNumber string) (bool, error) {
	log.Info().Msg("implement CheckBlacklist")
	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginBlackListOTP)

	count, err := collection.CountDocuments(ctx, bson.M{"phone_number": phoneNumber})
	return count > 0, err
}

func (or *otpRepository) MarkOTPAsObsolete(ctx context.Context, phoneNumber, code string) error {
	log.Info().Msg("implement MarkOTPAsObsolete")

	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginLogOTP)
	_, err := collection.UpdateOne(ctx, bson.M{"phone_number": phoneNumber, "code": code}, bson.M{"$set": bson.M{"obsolete": true}})

	return err
}

func (or *otpRepository) FindFailedAttempt(ctx context.Context, phoneNumber, code string) (*domain.FailedAttempts, error) {
	log.Info().Msg("implement FindFailedAttempt")

	var failedAttempts domain.FailedAttempts
	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginFailedAttemptsOTP)

	err := collection.FindOne(ctx, bson.M{"phone_number": phoneNumber, "code": code}).Decode(&failedAttempts)

	return &failedAttempts, err
}

func (or *otpRepository) LogFailedAttempt(ctx context.Context, clientID, phoneNumber, code string) (int, error) {
	log.Info().Msg("implement LogFailedAttempt")

	collection := or.mongo.Database(enums.MongodbDatabase).Collection(enums.LoginFailedAttemptsOTP)
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	filter := bson.M{"phone_number": phoneNumber, "code": code}
	update := bson.M{"$inc": bson.M{"attempts": 1}, "$set": bson.M{"client_id": clientID, "code": code, "failed_at": time.Now()}}

	result := domain.FailedAttempts{}
	err := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&result)

	return result.Attempts, err
}
