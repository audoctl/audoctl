package enum

type EventSourceType string

const (
	EventSourceTypeAPI    EventSourceType = "api"
	EventSourceTypeSystem EventSourceType = "system"
	EventSourceTypeAI     EventSourceType = "ai"
	EventSourceTypeUser   EventSourceType = "user"
)

func (e EventSourceType) String() string {
	return string(e)
}
