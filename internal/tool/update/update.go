package update

import (
	"io"

	"github.com/frioux/leatherman/internal/selfupdate"
)

/*
Update checks to see if there's an update from github and installs it if there
is.  If LM_GH_TOKEN is set to a personal access token this can be called more
frequently without exhausting github api limits.

Command: update
*/
func Update([]string, io.Reader) error {
	selfupdate.MaybeUpdate()

	return nil
}
