package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StopFn func(context.Context) error

type Handler struct {
	sigCh    chan os.Signal
	deadline time.Duration
}

func New(ctxDeadline time.Duration) Handler {
	return Handler{
		sigCh:    make(chan os.Signal, 1),
		deadline: ctxDeadline,
	}
}

func (h Handler) Handle(ctx context.Context, stopFn StopFn) {
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(h.sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until we receive our signal.
	<-h.sigCh

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(ctx, h.deadline)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if stopFn(ctx) != nil {
		os.Exit(1)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other api
	// to finalize based on context cancellation.
	fmt.Println("gracefully shutdown")
	os.Exit(0)
}
