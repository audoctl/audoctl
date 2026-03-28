package docs

import "github.com/audoctl/audoctl/internal/model/entity"

// swagger:route POST /api/v1/sessions sessions CreateSession
// Create a new session
//
// Responses:
//
//	201: SessionResponse
//	400: ErrorResponse
//	500: ErrorResponse
//
// swagger:parameters CreateSession
type CreateSessionRequest struct {
	// in: body
	Body struct {
		ActorID string `json:"actor_id"`
		Type    string `json:"type"`
	} `json:"body"`
}

// swagger:response SessionResponse
type SessionResponse struct {
	// in: body
	Body struct {
		Session entity.Session `json:"session"`
	} `json:"body"`
}
