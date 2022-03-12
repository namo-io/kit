package cmd

import (
	"strings"

	"github.com/namo-io/kit/pkg/buildinfo"
	"github.com/spf13/cobra"
)

// NewDefaultRootCmd create default root command line object
func NewDefaultRootCmd() *cobra.Command {
	version := buildinfo.GetVersion()

	cmd := &cobra.Command{
		Version: version,
	}

	if strings.Contains(".", version) {
		cmd.SetVersionTemplate(`{{ printf "v%s\n" .Version }}`)
	}

	return cmd
}
