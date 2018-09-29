package dumpmozlz4

import (
	"fmt"
	"io"
	"os"

	"github.com/frioux/mozlz4"
	"github.com/pkg/errors"
)

const magicHeader = "mozLz40\x00"

// DumpMozLZ4 writes the uncompressed mozlz4 file from the first argument to stdout
func DumpMozLZ4(args []string, _ io.Reader) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: %s session.jsonlz4", args[0])
	}
	file, err := os.Open(args[1])
	if err != nil {
		return errors.Wrap(err, "Couldn't open")
	}

	r, err := mozlz4.NewReader(file)
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		return errors.Wrap(err, "Couldn't copy")
	}

	return nil
}
