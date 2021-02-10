package update

import (
	"io"

	"github.com/frioux/leatherman/internal/selfupdate"
)

func Update([]string, io.Reader) error {
	selfupdate.MaybeUpdate()

	return nil
}
