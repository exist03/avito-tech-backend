package repository

import (
	"avito-tech-backend/config"
	"avito-tech-backend/domain"
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

func (r *PsqlRepo) Get(ctx context.Context, userId int) ([]domain.Segment, error) {
	err := r.checkUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, "SELECT segment.id, segment.name from segment join accordance on segment.id = accordance.segment_id where accordance.user_id=$1", userId)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoContent
		}
		return nil, err
	}

	var res []domain.Segment
	for rows.Next() {
		var segmentId int
		var segmentName string
		err = rows.Scan(&segmentId, &segmentName)
		if err != nil {
			return nil, err
		}
		segment := domain.Segment{Id: segmentId, Name: segmentName}
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
	return filename, nil
}

func (r *PsqlRepo) Create(ctx context.Context, segment domain.Segment) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO segment VALUES ($1, $2)", segment.Id, segment.Name)
	if err != nil {
		r.logger.Debug().Err(err).Msg("create error")
		return err
	}
	id, err := r.getPercentUsers(ctx, segment.Percent)
	if err != nil {
		return err
	}
	for i := 0; i < len(id); i++ {
		r.connectSegment(ctx, id[i], segment.Id, segment.TTL)
	}
	return nil
}
func (r *PsqlRepo) getPercentUsers(ctx context.Context, percent float64) ([]int, error) {
	var res []int
	rows, err := r.pool.Query(ctx, "SELECT id FROM users ORDER BY random() LIMIT (SELECT count(*)*$1/100 FROM users)", percent)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *PsqlRepo) Update(ctx context.Context, req domain.UpdateRequest) error {
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
		return
	}
	r.pool.Exec(ctx, "INSERT INTO history VALUES ($1, $2, $3, $4)", userId, segmentId, true, time.Now().Unix())
}
func (r *PsqlRepo) disconnectSegment(ctx context.Context, userId, segmentId int) {
	_, err := r.pool.Exec(ctx, "DELETE FROM accordance WHERE user_id=$1 AND segment_id=$2", userId, segmentId)
	if err != nil {
		r.logger.Info().Err(err).Msg("disconnect error")
		return
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

func (r *PsqlRepo) Checker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
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
				if err = rows.Scan(&userId, &segmentId); err != nil {
					continue
				}
				r.pool.Exec(context.Background(), "INSERT INTO history VALUES ($1, $2, $3, $4)", userId, segmentId, false, time.Now().Unix())
				r.pool.Exec(context.Background(), "DELETE FROM accordance WHERE user_id=$1 AND segment_id=$2", userId, segmentId)
			}
			rows.Close()
		}
	}
}
