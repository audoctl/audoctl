package enum

type EventType string

const (
	EventUserLogin EventType = "user.login"
	EventClick     EventType = "ui.click"
	EventAIStart   EventType = "ai.task.start"
	EventAIFinish  EventType = "ai.task.finish"
	EventAIFailure EventType = "ai.task.failure"
	EventAITask    EventType = "ai.task.task"
)

func (e EventType) String() string {
	return string(e)
}
