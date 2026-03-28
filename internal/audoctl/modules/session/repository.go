package session

import (
	"context"

	"github.com/audoctl/audoctl/internal/model/entity"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSession(ctx context.Context, session *entity.Session) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateSession(ctx context.Context, session *entity.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}
