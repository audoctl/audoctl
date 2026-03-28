package audoctl

import (
	"github.com/audoctl/audoctl/internal/audoctl/modules/event"
	"github.com/audoctl/audoctl/internal/audoctl/modules/session"
	"github.com/audoctl/audoctl/internal/audoctl/modules/timeline"
	"github.com/audoctl/audoctl/pkg/fiberserver"
)

type Handlers struct {
	Event    event.Handler
	Timeline timeline.Handler
	Session  session.Handler
}

func (hg Handlers) Handlers() []fiberserver.HandlerGroup {
	return []fiberserver.HandlerGroup{
		hg.Event,
		hg.Timeline,
		hg.Session,
	}
}
