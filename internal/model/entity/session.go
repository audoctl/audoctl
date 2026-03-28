package entity

import (
	"time"

	"github.com/audoctl/audoctl/internal/model/obj/id"
)

type Session struct {
	BaseEntity
	ActorID string
	Type    string
}

func NewSession(actorID string, sessionType string) *Session {
	return &Session{
		BaseEntity: BaseEntity{
			Id:        id.New(),
			CreatedAt: time.Now(),
		},
		ActorID: actorID,
		Type:    sessionType,
	}
}
