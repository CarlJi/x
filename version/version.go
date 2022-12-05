package version

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
)

const (
	unknownProperty = "N/A"
)

// Information of versioning
var (
	GoVersion     = unknownProperty
	GitCommit     = unknownProperty
	GitCommitDate = unknownProperty
	GitTag        = unknownProperty
	BuildDate     = unknownProperty
	Platform      = unknownProperty
)

func init() {
	// 支持 version、--version 两种模式输出version信息
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		Version()
		os.Exit(0)
	}

	// TODO: expose version information via exporter as needed
}

var once sync.Once

// Version prints the information of versioning
func Version() {
	once.Do(func() {
		collectFromBuildInfo()
		collectFromRuntime()
	})

	format := "%s:\t%s\n"
	xprintf(format, "Go version", GoVersion)
	xprintf(format, "Git commit", GitCommit)
	xprintf(format, "Commit date", GitCommitDate)
	xprintf(format, "Built date", BuildDate)
	xprintf(format, "Git tag", GitTag)
	xprintf(format, "OS/Arch", Platform)
}

// collectFromBuildInfo tries to set the build information embedded in the running binary via Go module.
// It doesn't override data if were already set by Go -ldflags.
func collectFromBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			if GitCommit == unknownProperty && kv.Value != "" {
				GitCommit = kv.Value
			}
		case "vcs.time":
			if GitCommitDate == unknownProperty && kv.Value != "" {
				GitCommitDate = kv.Value
			}
		}
	}
}

// collectFromRuntime tries to set the build information embedded in the running binary via go runtime.
// It doesn't override data if were already set by Go -ldflags.
func collectFromRuntime() {
	if GoVersion == unknownProperty {
		GoVersion = runtime.Version()
	}

	if Platform == unknownProperty {
		Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

// xprintf prints a message to standard output.
func xprintf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
