package cmd

import (
	"strings"

	"github.com/namo-io/kit/pkg/buildinfo"
	"github.com/namo-io/kit/pkg/log"
)

type CmdOption func(*Cmd)

func WithVersion() CmdOption {
	version := buildinfo.GetVersion()

	return func(c *Cmd) {
		c.c.Version = version

		// if version is semantic format(v1.1.1), setup template
		if strings.Contains(".", version) {
			c.c.SetVersionTemplate(`{{ printf "v%s\n" .Version }}`)
		}

		log.Tracef("this service version is '%v'", version)
	}
}

func WithConfiguration() CmdOption {
	return func(c *Cmd) {
		c.c.Flags().StringVarP(&c.configFilePath, "config", "c", "./config/config.yaml", "config file path (default: ./config/config.yaml)")
	}
}
