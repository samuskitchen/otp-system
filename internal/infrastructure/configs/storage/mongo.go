package storage

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"otp-system/pkg/kit/enums"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type key string

var (
	once     sync.Once
	instance *mongo.Client
)

const (
	dbMongo key = key(enums.BDMongo)
)

func ConnInstance() (client *mongo.Client) {
	once.Do(func() {
		instance = getConnections(dbMongo, os.Getenv("MONGO_CRED_URI"), enums.MongodbDatabase)
	})

	return instance
}

func getConnections(key key, mongoCredURI, databaseName string) *mongo.Client {
	ctx := context.WithValue(context.Background(), key, enums.DBConnection)

	connectTimeout := time.Duration(enums.MongodbMaxConnectionTimeOut) * time.Millisecond
	socketTimeout := time.Duration(enums.MongodbSocketTimeout) * time.Millisecond
	maxConnIdleTime := time.Duration(enums.MongodbMaxConnectionIdleTime) * time.Millisecond
	minPoolSize := enums.MongodbMinConnectionsPerHost
	maxPoolSize := enums.MongodbMaxConnectionsPerHost
	dbTimeout := time.Duration(enums.MongodbMaxDatabaseTimeOut) * time.Millisecond

	clientOptions := options.Client()
	clientOptions.ConnectTimeout = &connectTimeout
	clientOptions.SocketTimeout = &socketTimeout
	clientOptions.MaxConnIdleTime = &maxConnIdleTime
	clientOptions.MaxPoolSize = &maxPoolSize
	clientOptions.MinPoolSize = &minPoolSize

	ctxTimeout, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	client, err := mongo.Connect(ctxTimeout, clientOptions.ApplyURI(mongoCredURI))
	if err != nil {
		panic(fmt.Sprintf(enums.MongoErrorConfiguration, err.Error()))
	}

	if err = client.Ping(ctxTimeout, readpref.Primary()); err != nil {
		panic(fmt.Sprintf(enums.MongoErrorConnection, err.Error()))
	}

	client.Database(databaseName)

	log.Info().Msg("Database successfully read connected")
	return client
}
