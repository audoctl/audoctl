package configs

type Log struct {
	App   App    `env:"APP" default:"[NO-APP-NAME]"`
	Level string `env:"LEVEL" default:"info"`
	Env   string `env:"ENV" default:"dev"`
}

type App string

const (
	Audoctl App = "[AUDOCTL]"
)

const (
	AudoctlAppName = "Audoctl"
)

func (t App) IsValid() bool {
	types := map[App]struct{}{
		Audoctl: {},
	}

	_, ok := types[t]
	return ok
}
