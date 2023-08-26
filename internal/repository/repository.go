package repository

import (
	"avito-tech-backend/config"
	"avito-tech-backend/domain"
	"avito-tech-backend/internal"
	"avito-tech-backend/pkg/logger"
	"avito-tech-backend/pkg/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"os"
	"time"
)

const filename = "file.csv"

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

func (r *PsqlRepo) Get(ctx context.Context, userId int) ([]internal.Segment, error) {
	err := r.checkUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	res := make([]internal.Segment, 0)
	rows, err := r.pool.Query(ctx, "SELECT segment.id, segment.name from segment join accordance on segment.id = accordance.segment_id where accordance.user_id=$1", userId)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoContent
		}
		return nil, err
	}
	for rows.Next() {
		var segmentId int
		var segmentName string
		err = rows.Scan(&segmentId, &segmentName)
		if err != nil {
			return nil, err
		}
		segment := internal.Segment{Id: segmentId, Name: segmentName}
		res = append(res, segment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *PsqlRepo) GetHistory(ctx context.Context, timeBegin, timeEnd int64, userId int) (string, error) {
	err := r.checkUser(ctx, userId)
	if err != nil {
		return "", err
	}
	rows, err := r.pool.Query(ctx, "SELECT segment_id, type, time FROM history where time BETWEEN $1 AND $2 AND user_id=$3", timeBegin, timeEnd, userId)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrNoContent
		}
		return "", err
	}
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	for rows.Next() {
		var (
			segmentId int
			t         bool
			unixTime  int64
		)

		if err = rows.Scan(&segmentId, &t, &unixTime); err != nil {
			return "", err
		}
		fmt.Fprintln(file, userId, ";", segmentId, ";", t, ";", unixTime)
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	return filename, nil // remove hardcode
}

func (r *PsqlRepo) Create(ctx context.Context, segment internal.Segment) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO segment VALUES ($1, $2)", segment.Id, segment.Name)
	if err != nil {
		r.logger.Debug().Err(err).Msg("create error")
		return err
	}
	return nil
}

func (r *PsqlRepo) Update(ctx context.Context, req internal.UpdateRequest) error {
	err := r.checkUser(ctx, req.UserId)
	if err != nil {
		return err
	}
	for _, v := range req.SegmentsAdd {
		r.connectSegment(ctx, req.UserId, v.Id, v.TTL)
	}
	for _, v := range req.SegmentsDel {
		r.disconnectSegment(ctx, req.UserId, v.Id)
	}
	return nil
}
func (r *PsqlRepo) connectSegment(ctx context.Context, userId, segmentId int, ttl time.Time) {
	_, err := r.pool.Exec(ctx, "INSERT INTO accordance VALUES ($1, $2, $3)", userId, segmentId, ttl.Unix())
	if err != nil {
		r.logger.Info().Err(err).Msg("connect error")
	}
	r.pool.Exec(ctx, "INSERT INTO history VALUES ($1, $2, $3, $4)", userId, segmentId, true, time.Now().Unix())
}
func (r *PsqlRepo) disconnectSegment(ctx context.Context, userId, segmentId int) {
	_, err := r.pool.Exec(ctx, "DELETE FROM accordance WHERE user_id=$1 AND segment_id=$2", userId, segmentId)
	if err != nil {
		r.logger.Info().Err(err).Msg("disconnect error")
	}
	r.pool.Exec(ctx, "INSERT INTO history VALUES ($1, $2, $3, $4)", userId, segmentId, false, time.Now().Unix())
}

func (r *PsqlRepo) Delete(ctx context.Context, segmentId int) error {
	res, err := r.pool.Exec(ctx, "DELETE FROM segment WHERE id=$1", segmentId)
	if err != nil {
		r.logger.Debug().Err(err).Msg("delete error")
		return err
	}
	if res.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PsqlRepo) checkUser(ctx context.Context, userId int) error {
	res, err := r.pool.Exec(ctx, "SELECT * FROM users WHERE id=$1", userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PsqlRepo) Checker() {
	for {
		time.Sleep(45 * time.Second)
		rows, err := r.pool.Query(context.Background(), "SELECT user_id, segment_id FROM accordance WHERE expires<$1", time.Now().Unix())
		if err != nil {
			continue
		}
		for rows.Next() {
			var (
				userId    int
				segmentId int
			)
			rows.Scan(&userId, &segmentId)
			r.pool.Exec(context.Background(), "INSERT INTO history VALUES ($1, $2, $3, $4)", userId, segmentId, false, time.Now().Unix())
			r.pool.Exec(context.Background(), "DELETE FROM accordance WHERE user_id=$1 AND segment_id=$2", userId, segmentId)
		}
	}
}
