package timeline

import (
	"context"

	"github.com/audoctl/audoctl/internal/model/entity"
	"gorm.io/gorm"
)

type Repository interface {
	GetEventsBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*entity.Event, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetEventsBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*entity.Event, error) {
	var events []*entity.Event
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
