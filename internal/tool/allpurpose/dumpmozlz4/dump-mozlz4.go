package dumpmozlz4

import (
	"fmt"
	"io"
	"os"

	"github.com/frioux/leatherman/pkg/mozlz4"
)

func Run(args []string, _ io.Reader) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: %s session.jsonlz4", args[0])
	}
	file, err := os.Open(args[1])
	if err != nil {
		return fmt.Errorf("Couldn't open: %w", err)
	}

	r, err := mozlz4.NewReader(file)
	if err != nil {
		return fmt.Errorf("mozlz4.NewReader: %w", err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		return fmt.Errorf("Couldn't copy: %w", err)
	}

	return nil
}
