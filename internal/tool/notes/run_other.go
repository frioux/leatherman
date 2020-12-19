// +build !linux

package notes

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func Serve(args []string, _ io.Reader) error {
	fmt.Fprintf(os.Stderr, "zine not supported on %s/%s due to lacking support in modernc.org/sqlite\n", runtime.GOOS, runtime.GOARCH)

	return nil
}
