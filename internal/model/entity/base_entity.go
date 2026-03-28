package entity

import (
	"time"

	"github.com/audoctl/audoctl/internal/model/obj/id"
	"github.com/audoctl/audoctl/internal/shared/helper"
)

type BaseEntity struct {
	Id        id.ID     `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type TimestampedEntity struct {
	BaseEntity
	DeletedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-" gorm:"updated_at"`
}

func NewTimestampedEntity() TimestampedEntity {
	return TimestampedEntity{
		BaseEntity: BaseEntity{
			Id:        id.New(),
			CreatedAt: time.Now(),
		},
	}
}

func (b *TimestampedEntity) Delete() {
	b.DeletedAt = helper.Ptr(time.Now())
}

type AuditableEntity struct {
	TimestampedEntity
	CreatedBy id.ID  `json:"createdBy,omitempty"`
	UpdatedBy *id.ID `json:"-"`
}

func NewAuditableEntity(userID id.ID) AuditableEntity {
	return AuditableEntity{
		TimestampedEntity: NewTimestampedEntity(),
		CreatedBy:         userID,
	}
}

func (a *AuditableEntity) Delete(updatedBy id.ID) {
	a.UpdatedBy = &updatedBy
	a.TimestampedEntity.Delete()
}

func (b *BaseEntity) IsNil() bool {
	return b == nil || b.Id.IsInvalid()
}
