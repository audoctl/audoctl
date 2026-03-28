package event

import (
	"context"

	"github.com/audoctl/audoctl/internal/model/entity"
	"github.com/audoctl/audoctl/internal/model/enum"
)

type Service interface {
	CreateEvent(ctx context.Context, req EventCreateRequest) (*entity.Event, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) CreateEvent(ctx context.Context, req EventCreateRequest) (*entity.Event, error) {
	event := entity.NewEvent(
		req.SessionID,
		enum.EventType(req.Type),
		enum.EventSourceType(req.Source),
		req.ActorID,
		req.Data,
	)
	if err := s.repository.CreateEvent(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}
