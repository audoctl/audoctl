//go:build wireinject
// +build wireinject

//go:generate wire

package audoctl

import (
	"github.com/audoctl/audoctl/configs"
	"github.com/audoctl/audoctl/internal/audoctl/modules/event"
	"github.com/audoctl/audoctl/internal/audoctl/modules/session"
	"github.com/audoctl/audoctl/internal/audoctl/modules/timeline"
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func InitHandlers(
	db *gorm.DB,
	cfg *configs.Config,
	logger zerolog.Logger,
) Handlers {
	wire.Build(
		event.NewRepository,
		event.NewService,
		event.NewHandler,

		timeline.NewRepository,
		timeline.NewService,
		timeline.NewHandler,

		session.NewRepository,
		session.NewService,
		session.NewHandler,

		wire.Struct(new(Handlers), "*"),
	)
	return Handlers{}
}
