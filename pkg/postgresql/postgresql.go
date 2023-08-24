package postgresql

import (
	"avito-tech-backend/config"
	"avito-tech-backend/pkg/logger"
	"avito-tech-backend/pkg/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewClient(ctx context.Context, config config.PsqlStorage, maxAttempts int) (pool *pgxpool.Pool, err error) {
	log := logger.GetLogger()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	log.Info().Msg(dsn)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			log.Err(err)
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
