package docs

import "github.com/audoctl/audoctl/internal/audoctl/modules/timeline"

// swagger:route GET /api/v1/timeline/{session_id} timeline GetEventsBySessionID
// Get events by session ID
//
// Responses:
//
//	200: EventsBySessionIDResponse
//	400: ErrorResponse
//	500: ErrorResponse
//
// swagger:parameters GetEventsBySessionID
type GetEventsBySessionID struct {
	// in: path
	SessionID string `json:"session_id"`
	// in: query
	Limit int `json:"limit" default:"10"`
	// in: query
	Offset int `json:"offset" default:"0"`
}

// swagger:response EventsBySessionIDResponse
type EventsBySessionIDResponse struct {
	// in: body
	Body struct {
		Events []*timeline.EventResponse `json:"events"`
	} `json:"body"`
}
