package cmd

import (
	"github.com/namo-io/kit/pkg/buildinfo"
	"github.com/spf13/cobra"
)

// NewDefaultRootCmd create default root command line object
func NewDefaultRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Version: buildinfo.GetVersion(),
	}
	cmd.SetVersionTemplate(`{{ printf "%s\n" .Version }}`)

	return cmd
}
