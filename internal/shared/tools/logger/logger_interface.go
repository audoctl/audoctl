package logger

type Logger interface {
	Info(msg string, args map[string]any)
	Error(msg string, args map[string]any)
	Warn(msg string, args map[string]any)
	Fatal(msg string, args map[string]any)
	Debug(msg string, args map[string]any)
	GetLevel() string
}
