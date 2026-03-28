package session

// internal/module/event/dto.go
type CreateSessionRequest struct {
	ActorID string `json:"actor_id"`
	Type    string `json:"type"`
}
