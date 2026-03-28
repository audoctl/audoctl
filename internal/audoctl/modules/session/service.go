package session

import (
	"context"

	"github.com/audoctl/audoctl/internal/model/entity"
)

type Service interface {
	CreateSession(ctx context.Context, req CreateSessionRequest) (*entity.Session, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) CreateSession(ctx context.Context, req CreateSessionRequest) (*entity.Session, error) {
	session := entity.NewSession(
		req.ActorID,
		req.Type,
	)
	if err := s.repository.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}
