package timeline

import (
	"time"

	"github.com/audoctl/audoctl/internal/model/enum"
	"github.com/audoctl/audoctl/internal/model/obj/id"
	"github.com/audoctl/audoctl/pkg/errs"
)

type GetEventsRequest struct {
	SessionID string `json:"session_id" uri:"session_id"`
	Limit     int    `json:"limit" query:"limit" default:"10"`
	Offset    int    `json:"offset" query:"offset" default:"0"`
}

func (r *GetEventsRequest) Validate() error {
	if r.SessionID == "" {
		return errs.BadRequest("Session ID is required")
	}
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Offset < 0 {
		r.Offset = 0
	}
	return nil
}

type EventResponse struct {
	Id        id.ID                  `json:"id"`
	SessionID string                 `json:"session_id"`
	Type      enum.EventType         `json:"type"`
	Source    enum.EventSourceType   `json:"source"`
	ActorID   string                 `json:"actor_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}
