package configs

type Application struct {
	Name    string `env:"NAME,default=Audoctl"`
	Version string `env:"VERSION,default=undefined"`
}
