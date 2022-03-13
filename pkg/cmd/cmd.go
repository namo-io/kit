package cmd

import (
	"context"
	"fmt"

	"github.com/namo-io/kit/pkg/buildinfo"
	"github.com/namo-io/kit/pkg/log"
	"github.com/namo-io/kit/pkg/metric"
	"github.com/namo-io/kit/pkg/trace"
	"github.com/namo-io/kit/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Cmd struct {
	c cobra.Command

	configFilePath string
}

type CmdRun func(context.Context, *viper.Viper) error

func New(opts ...CmdOption) *Cmd {
	cmd := &Cmd{}

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func (c *Cmd) SetRun(fn CmdRun) *Cmd {
	c.c.Run = func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		v := viper.New()

		// if used WithConfiguration option
		if len(c.configFilePath) != 0 {
			log.WithField("config.file.path", c.configFilePath).Infof("configuration loading...")

			v.SetConfigFile(c.configFilePath)
			if err := v.ReadInConfig(); err != nil {
				log.Fatal(err)
			}
		}

		// global setup
		if err := setupLogging(v); err != nil {
			log.Fatal(err)
		}

		if err := setupTracing(v); err != nil {
			log.Fatal(err)
		}

		if err := setupMetric(v); err != nil {
			log.Fatal(err)
		}

		// origin runner
		if err := fn(ctx, v); err != nil {
			log.Fatal(err)
		}
	}

	return c
}

func (c *Cmd) Execute() error {
	return c.c.Execute()
}

func setupLogging(v *viper.Viper) error {
	v.SetDefault("logging.verbose", true)

	verbose := v.GetBool("logging.verbose")
	serviceName := v.GetString("logging.service.name")
	flog := log.WithField("logging.verbose", fmt.Sprintf("%v", verbose))

	log.SetVerbose(verbose)
	log.SetFields(map[string]string{
		"service.id":      util.GetServiceId(),
		"service.version": buildinfo.GetVersion(),
		"service.host":    util.GetHostname(),
	})

	if len(serviceName) != 0 {
		flog = flog.WithField("logging.service.name", serviceName)
		log.SetField("service.name", serviceName)
	}

	flog.Trace("global logging setup")
	return nil
}

func setupTracing(v *viper.Viper) error {
	tracingEnabled := v.GetBool("tracing.enabled")
	jeagerEndpoint := v.GetString("tracing.jeager.endpoint")
	serviceName := v.GetString("tracing.service.name")

	if tracingEnabled {
		if len(jeagerEndpoint) != 0 {
			if len(serviceName) == 0 {
				return fmt.Errorf("tracing service name parameter is required, but it is empty")
			}

			err := trace.SetJeagerTraceProvider(
				serviceName,
				util.GetServiceId(),
				buildinfo.GetVersion(),
				jeagerEndpoint,
			)
			if err != nil {
				return err
			}

			log.WithField("tracing.service.name", serviceName).
				WithField("tracing.jeager.endpoint", jeagerEndpoint).
				Trace("global tracing setup")
		}
	}

	return nil
}

func setupMetric(v *viper.Viper) error {
	metricEnabled := v.GetBool("metric.enabled")
	metricPort := v.GetInt("metric.port")
	serviceName := v.GetString("metric.service.name")

	if metricEnabled {
		err := metric.SetPrometheusMetricProvider(
			serviceName,
			util.GetServiceId(),
			buildinfo.GetVersion(),
			metricPort,
		)

		if err != nil {
			return err
		}

		log.WithField("metric.port", fmt.Sprintf("%v", metricPort)).
			Trace("global metric setup")
	}
	return nil
}
