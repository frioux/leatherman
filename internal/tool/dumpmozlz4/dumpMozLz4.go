package dumpmozlz4

import (
	"io"
	"os"

	"github.com/frioux/leatherman/pkg/mozlz4"
	"golang.org/x/xerrors"
)

// Run writes the uncompressed mozlz4 file from the first argument to stdout
func Run(args []string, _ io.Reader) error {
	if len(args) != 2 {
		return xerrors.Errorf("Usage: %s session.jsonlz4", args[0])
	}
	file, err := os.Open(args[1])
	if err != nil {
		return xerrors.Errorf("Couldn't open: %w", err)
	}

	r, err := mozlz4.NewReader(file)
	if err != nil {
		return xerrors.Errorf("mozlz4.NewReader: %w", err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		return xerrors.Errorf("Couldn't copy: %w", err)
	}

	return nil
}
