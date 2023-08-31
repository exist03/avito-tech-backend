package service

import (
	"avito-tech-backend/domain"
	"context"
)

//go:generate mockery --name Repository --with-expecter=true
type Repository interface {
	Get(ctx context.Context, userId int) ([]domain.Segment, error)
	Create(ctx context.Context, segment domain.Segment) error
	Delete(ctx context.Context, segmentId int) error
	Update(ctx context.Context, req domain.UpdateRequest) error
	GetHistory(ctx context.Context, timeBegin, timeEnd int64, userId int) (string, error)
}
type Service struct {
	repo Repository
}

func New(repository Repository) *Service {
	return &Service{repo: repository}
}

func (s *Service) Get(userId int) ([]domain.Segment, error) {
	list, err := s.repo.Get(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) GetHistory(timeBegin, timeEnd int64, userId int) (string, error) {
	if timeEnd < timeBegin {
		return "", domain.ErrInvalidArgument
	}
	path, err := s.repo.GetHistory(context.Background(), timeBegin, timeEnd, userId)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s *Service) Create(segment domain.Segment) error {
	return s.repo.Create(context.Background(), segment)
}

func (s *Service) Update(req domain.UpdateRequest) error {
	return s.repo.Update(context.Background(), req)
}

func (s *Service) Delete(segmentId int) error {
	return s.repo.Delete(context.Background(), segmentId)
}
