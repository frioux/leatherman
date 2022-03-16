package version

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	BS = bi.Settings[:0]
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" {
			Version = s.Value
			continue
		}
		if s.Key == "vcs.time" {
			When = s.Value
			continue
		}
		BS = append(BS, s)
	}

	Deps = bi.Deps
}

// Version is the git version that produced this binary.
var Version string

// When is the datestamp that produced this binary.
var When string

var BS []debug.BuildSetting

var Deps []*debug.Module

func Render(w io.Writer) {
	fmt.Fprintf(w, "Leatherman built from %s on %s by with %s\n",
		Version, When, runtime.Version())

	fmt.Fprintln(w, "Build Settings:")
	for _, s := range BS {
		fmt.Fprintf(w, "\t%s=%s\n", s.Key, s.Value)
	}

	fmt.Fprintln(w, "\nDeps:")
	for _, dep := range Deps {
		fmt.Fprintf(w, "\t%s@%s (%s)\n", dep.Path, dep.Version, dep.Sum)
		if dep.Replace != nil {
			r := dep.Replace
			fmt.Fprintf(w, "   replaced by %s@%s (%s)\n", r.Path, r.Version, r.Sum)
		}
	}
}
