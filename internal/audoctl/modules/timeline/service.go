package timeline

import (
	"context"
)

type Service interface {
	GetEventsBySessionID(ctx context.Context, req GetEventsRequest) ([]*EventResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetEventsBySessionID(ctx context.Context, req GetEventsRequest) ([]*EventResponse, error) {
	events, err := s.repository.GetEventsBySessionID(ctx, req.SessionID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	response := make([]*EventResponse, len(events))
	for i, event := range events {
		response[i] = &EventResponse{
			Id:        event.Id,
			SessionID: event.SessionID,
			Type:      event.Type,
			Source:    event.Source,
			ActorID:   event.ActorID,
			Data:      event.Data,
			CreatedAt: event.CreatedAt,
		}
	}
	return response, nil
}
