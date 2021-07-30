package cmds

import (
	"fmt"

	"github.com/namo-io/kit/pkg/version"
	"github.com/spf13/cobra"
)

func Version() *cobra.Command {
	var more bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version",
		Long:  "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !more {
				fmt.Println(version.String())
				return nil
			}
			fmt.Println(version.Info().String())

			return nil
		},
	}
	cmd.PersistentFlags().BoolVarP(&more, "more", "m", false, "show more detail information")

	return cmd
}
