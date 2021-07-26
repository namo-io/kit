package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"

	"github.com/ryanuber/columnize"
)

// Exported variables are expected to be set via -ldflags -X options at build-time
// ex) go build ./cmd -ldflags -X Buildtime=12:12:12 ...
var (
	// Buildtime is the time at created build artifact
	Buildtime string

	// CommitHash is the git commit hash
	CommitHash string

	// Tag is the git recent tag
	Tag string

	v version
)

type version struct {
	Buildtime    string            `json:"buildtime"`
	CommitHash   string            `json:"commithash"`
	Tag          string            `json:"tag"`
	Goversion    string            `json:"goversion"`
	OS           string            `json:"os"`
	Architecture string            `json:"architecture"`
	Modules      map[string]string `json:"modules"`
}

func init() {
	v = version{
		Tag:          Tag,
		Buildtime:    Buildtime,
		CommitHash:   CommitHash,
		Goversion:    runtime.Version(),
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Modules:      getModules(),
	}
}

func getModules() (modules map[string]string) {
	modules = make(map[string]string)
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, req := range buildInfo.Deps {
			path := req.Path
			modules[path] = req.Version
		}
	}

	return
}

func String() string {
	if v.Tag == "" {
		return fmt.Sprintf("%v", v.CommitHash)
	}

	if v.CommitHash == "" {
		return fmt.Sprintf("%v", v.Tag)
	}

	return "Unknown"
}

func Info() version {
	return v
}

func (v version) String() string {
	cfg := columnize.DefaultConfig()
	cfg.NoTrim = true
	cfg.Prefix = " "

	moduleVersion := []string{}
	for path, version := range v.Modules {
		moduleVersion = append(moduleVersion, fmt.Sprintf("·%s|%s", path, version))
	}
	sort.Strings(moduleVersion)

	return columnize.Format([]string{
		fmt.Sprintf("Version:|%s", String()),
		fmt.Sprintf("Build Time:|%s", v.Buildtime),
		fmt.Sprintf("Go Version:|%s", v.Goversion),
		fmt.Sprintf("Git Commit:|%s", v.CommitHash),
		fmt.Sprintf("OS/Arch:|%s/%s", v.OS, v.Architecture),
		"Modules:|",
	}, cfg) + fmt.Sprintf("\r\n%s", columnize.Format(moduleVersion, cfg))
}
