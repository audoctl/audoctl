package audoctl

import (
	"context"
	"fmt"
	"log"

	"github.com/audoctl/audoctl/configs"
	"github.com/audoctl/audoctl/pkg/graceful"
	"github.com/spf13/cobra"
)

var CtlCmd = &cobra.Command{
	Use:   "audoctl",
	Short: "Audit and execution control for AI systems",
	Run:   startCtl,
}

func startCtl(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cfg, err := configs.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	cfg.Log.App = configs.Audoctl
	cfg.Application.Name = configs.AudoctlAppName
	stopFn := StartCtl(ctx, cfg)

	graceful.New(configs.DefaultContextDeadline).Handle(ctx, stopFn)
	fmt.Println("audoctl started")
	fmt.Println(cfg)
}
