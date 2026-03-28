package event

type EventCreateRequest struct {
	SessionID string                 `json:"session_id" validate:"required"`
	Type      string                 `json:"type" validate:"required"`
	Source    string                 `json:"source" validate:"required,oneof=api ai system user"`
	ActorID   string                 `json:"actor_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
