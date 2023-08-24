package service

import (
	"avito-tech-backend/internal"
	"context"
	"encoding/json"
)

type Repository interface {
	Get(ctx context.Context, userId int) ([]string, error)
	Create(ctx context.Context, segment internal.Segment) error
	Delete(ctx context.Context, segmentId int) error
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
	res := internal.SegmentList{Segments: list}
	response, _ := json.Marshal(res)
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
