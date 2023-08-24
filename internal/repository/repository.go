package repository

import (
	"avito-tech-backend/config"
	"avito-tech-backend/internal"
	"avito-tech-backend/pkg/logger"
	"avito-tech-backend/pkg/postgresql"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PsqlRepo struct {
	pool   *pgxpool.Pool
	logger zerolog.Logger
}

func NewPsql(ctx context.Context, config config.PsqlStorage) *PsqlRepo {
	log := logger.GetLogger()
	pool, err := postgresql.NewClient(ctx, config, 3)
	if err != nil {
		log.Fatal().Err(err).Msg("Can`t create psql client")
	}
	return &PsqlRepo{pool: pool, logger: log}
}

func (r *PsqlRepo) Create(ctx context.Context, segment internal.Segment) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO segment VALUES ($1, $2)", segment.Id, segment.Name)
	if err != nil {
		r.logger.Debug().Err(err).Msg("create error")
		return err
	}
	return nil
}

func (r *PsqlRepo) Delete(ctx context.Context, segmentId int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM segment WHERE id=$1", segmentId)
	if err != nil {
		r.logger.Debug().Err(err).Msg("delete error")
		return err
	}
	return nil
}

func (r *PsqlRepo) Get(ctx context.Context, userId int) ([]internal.Segment, error) {
	res := make([]internal.Segment, 0)
	rows, err := r.pool.Query(ctx, "SELECT segment.id, segment.name from segment join accordance on segment.id = accordance.segment_id where accordance.user_id=$1", userId)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	for rows.Next() {
		var segmentId int
		var segmentName string
		err := rows.Scan(&segmentId, &segmentName)
		if err != nil {
			return nil, err
		}
		segment := internal.Segment{Id: segmentId, Name: segmentName}
		res = append(res, segment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

//func (r *PsqlRepo) Update(ctx context.Context) error {
//
//}
