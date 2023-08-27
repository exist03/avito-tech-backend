package service

import (
	"avito-tech-backend/domain"
	"avito-tech-backend/internal"
	"context"
	"encoding/json"
)

//go:generate mockery --name Repository
type Repository interface {
	Get(ctx context.Context, userId int) ([]internal.Segment, error)
	Create(ctx context.Context, segment internal.Segment) error
	Delete(ctx context.Context, segmentId int) error
	Update(ctx context.Context, req internal.UpdateRequest) error
	GetHistory(ctx context.Context, timeBegin, timeEnd int64, userId int) (string, error)
}
type Service struct {
	repo Repository
}

func New(repository Repository) *Service {
	return &Service{repo: repository}
}

func (s *Service) Get(userId int) ([]byte, error) {
	list, err := s.repo.Get(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	response, _ := json.Marshal(list)
	return response, nil
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

func (s *Service) Create(segment internal.Segment) error {
	err := s.repo.Create(context.Background(), segment)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(req internal.UpdateRequest) error {
	err := s.repo.Update(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(segmentId int) error {
	err := s.repo.Delete(context.Background(), segmentId)
	if err != nil {
		return err
	}
	return nil
}
