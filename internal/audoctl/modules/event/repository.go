package event

import (
	"context"

	"github.com/audoctl/audoctl/internal/model/entity"
	"gorm.io/gorm"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *entity.Event) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateEvent(ctx context.Context, event *entity.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}
