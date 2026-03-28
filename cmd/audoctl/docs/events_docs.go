package docs

import (
	"encoding/json"

	"github.com/audoctl/audoctl/internal/model/entity"
)

// swagger:route POST /api/v1/events events CreateEvent
// Create a new event
//
// Responses:
//
//	201: EventResponse
//	400: ErrorResponse
//	500: ErrorResponse
//
// swagger:parameters CreateEvent
type CreateEventRequest struct {
	// in: body
	Body struct {
		SessionID string          `json:"session_id"`
		Type      string          `json:"type"`
		Source    string          `json:"source"`
		ActorID   string          `json:"actor_id"`
		Data      json.RawMessage `json:"data"`
	} `json:"body"`
}

// swagger:response EventResponse
type EventResponse struct {
	// in: body
	Body struct {
		Event entity.Event `json:"event"`
	} `json:"body"`
}
