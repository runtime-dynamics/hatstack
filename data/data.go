package data

import (
	"context"
	"errors"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/rs/zerolog/log"
	"runtime-dynamics/config"
)

var (
	dataOnce   sync.Once
	dataClient *datastore.Client
	ctx        context.Context
)

func Cli() *datastore.Client {
	var err error
	dataOnce.Do(func() {
		ctx = context.Background()
		dataClient, err = datastore.NewClientWithDatabase(ctx, config.Get().GoogleProjectID, config.Get().DataStoreName)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create datastore client")
		}
	})
	return dataClient
}

func IsNotFound(err error) bool {
	return errors.Is(err, datastore.ErrNoSuchEntity)
}
