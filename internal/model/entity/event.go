package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/audoctl/audoctl/internal/model/enum"
	"github.com/audoctl/audoctl/internal/model/obj/id"
)

type Event struct {
	BaseEntity
	SessionID string               `json:"sessionId,omitempty"`
	Type      enum.EventType       `json:"type,omitempty"`
	Source    enum.EventSourceType `json:"source,omitempty"`
	ActorID   string               `json:"actorId,omitempty"`
	Data      EventData            `json:"data,omitempty"`
}

type EventData map[string]any

func (e *EventData) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), e)
}

func (e *EventData) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func NewEvent(sessionID string, eventType enum.EventType, source enum.EventSourceType, actorID string, data map[string]any) *Event {
	return &Event{
		BaseEntity: BaseEntity{
			Id:        id.New(),
			CreatedAt: time.Now(),
		},
		SessionID: sessionID,
		Type:      eventType,
		Source:    source,
		ActorID:   actorID,
		Data:      data,
	}
}
