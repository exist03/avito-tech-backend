package service

import (
	"avito-tech-backend/internal"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	Get(ctx context.Context, userId int) ([]internal.Segment, error)
	Create(ctx context.Context, segment internal.Segment) error
	Delete(ctx context.Context, segmentId int) error
	Update(ctx context.Context, req internal.UpdateRequest) error
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
		//TODO
	}
	response, _ := json.Marshal(list)
	return response, err
}

func (s *Service) Create(segment internal.Segment) error {
	err := s.repo.Create(context.Background(), segment)
	if err != nil {
		//TODO разделить ошибку сервера и ошибку ввода
		return err
	}
	return nil
}

func (s *Service) Delete(segmentId int) error {
	err := s.repo.Delete(context.Background(), segmentId)
	if err != nil {
		//TODO разделить ошибку сервера и ошибку ввода
		return err
	}
	return nil
}

func (s *Service) Update(req internal.UpdateRequest) error {
	err := s.repo.Update(context.Background(), req)
	if err != nil {
		log.Err(err)
	}
	return nil
}
