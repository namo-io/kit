package log

import (
	"context"
	"os"
	"runtime"

	"github.com/namo-io/kit/pkg/log/logger"
	"github.com/namo-io/kit/pkg/log/logger/typist"
	hooks "github.com/namo-io/kit/pkg/log/logger/typist/hookers"
)

var (
	_, fileName, _, _ = runtime.Caller(0)

	global = typist.New(
		typist.WithCllerIgnorePackageFile(fileName),
	)
)

const (
	ApplicationMetaKey = "app"
	ComponentMetaKey   = "component"
	HostnameMetaKey    = "host"
)

func init() {
	global.AddHooker(hooks.NewContextMapper())
}

type Config struct {
	// Verbose print include debug, trace messages
	Verbose bool

	// DisableColor print without colors
	DisableColor bool

	// Application log with field name "application"
	Application string

	ElasticSearchHookerConfig ElasticSearchHookerConfig
}

type ElasticSearchHookerConfig struct {
	Enable bool

	// Endpoint Elasticsearch hook endpoint (e.g. "localhost:9200")
	Endpoint string

	// Index Elasticsearch index name (e.g. "spc-cicd-api-logs")
	Index string

	// Elasticsearch Client Sniff (e.g. false)
	// docs: https://github.com/olivere/elastic/wiki/Sniffing
	Sniff bool
}

func AddConfig(cfg Config) error {
	hostname, err := os.Hostname()
	if err == nil {
		AddField(HostnameMetaKey, hostname)
	} else {
		AddField(HostnameMetaKey, "Unknown")
	}

	if !cfg.Verbose {
		AddOption(typist.WithVerboseDisable())
	}

	if len(cfg.Application) > 0 {
		AddField(ApplicationMetaKey, cfg.Application)
	}

	if cfg.DisableColor {
		formatter := typist.NewDefaultFormatter()
		formatter.IsUseColor = false
		AddOption(typist.WithFormatter(formatter))
	}

	esHookerCfg := cfg.ElasticSearchHookerConfig
	if esHookerCfg.Enable {
		es, err := hooks.NewElasticSearchSender(&hooks.ElasticSearchSenderConfig{
			Endpoints:       []string{esHookerCfg.Endpoint},
			Sniff:           esHookerCfg.Sniff,
			Index:           esHookerCfg.Index,
			IsCallStackFire: true,
		})
		if err != nil {
			return err
		}

		global.AddHooker(es)
	}

	return nil
}

// AddOption add options to global logger
func AddOption(opt typist.Options) {
	opt(global)
}

// AddHooker adding hooker to global logger
func AddHooker(hooker typist.Hooker) error {
	return global.AddHooker(hooker)
}

// AddField adding field to global logger
func AddField(key string, value interface{}) {
	global.AddField(key, value)
}

func Debug(args ...interface{}) {
	global.Debug(args...)
}
func Info(args ...interface{}) {
	global.Info(args...)
}
func Warn(args ...interface{}) {
	global.Warn(args...)
}
func Error(args ...interface{}) {
	global.Error(args...)
}
func Fatal(args ...interface{}) {
	global.Fatal(args...)
}
func Trace(args ...interface{}) {
	global.Trace(args...)
}

func Debugf(format string, args ...interface{}) {
	global.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	global.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	global.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	global.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	global.Fatalf(format, args...)
}
func Tracef(format string, args ...interface{}) {
	global.Tracef(format, args...)
}

func WithField(key string, value interface{}) logger.Logger {
	return global.WithField(key, value)
}
func WithFields(values map[string]interface{}) logger.Logger {
	return global.WithFields(values)
}
func WithContext(ctx context.Context) logger.Logger {
	return global.WithContext(ctx)
}
