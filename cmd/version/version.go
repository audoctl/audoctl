package version

import (
	"fmt"
	"runtime"

	"github.com/audoctl/audoctl/configs"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// Commit is set at build time via ldflags
	Commit = "unknown"
	// BuildTime is set at build time via ldflags
	BuildTime = "unknown"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  "Display detailed version information including build time, commit hash, and Go version",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("\n")
	fmt.Printf("╔═══════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                   AUDOCTL VERSION                         ║\n")
	fmt.Printf("╠═══════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Version:     %-44s ║\n", getVersion())
	fmt.Printf("║ Commit:      %-44s ║\n", Commit)
	fmt.Printf("║ Built:       %-44s ║\n", BuildTime)
	fmt.Printf("║ Go Version:  %-44s ║\n", runtime.Version())
	fmt.Printf("║ OS/Arch:     %-44s ║\n", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	fmt.Printf("╚═══════════════════════════════════════════════════════════╝\n")
	fmt.Printf("\n")
	fmt.Printf("  Audit and execution control for AI systems\n")
	fmt.Printf("  https://github.com/audoctl/audoctl\n")
	fmt.Printf("\n")
}

func getVersion() string {
	if Version == "dev" || Version == "" {
		return configs.Version
	}
	return Version
}
