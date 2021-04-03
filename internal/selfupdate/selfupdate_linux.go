package selfupdate

import (
	"fmt"
	"os"
	"syscall"
)

func isSameFile(path string) error {
	statExe, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("couldn't stat %s: %w", path, err)
	}

	statSelf, err := os.Stat("/proc/self/exe")
	if err != nil {
		return fmt.Errorf("couldn't stat /proc/self/exec: %w", err)
	}

	if statExe.Sys().(*syscall.Stat_t).Ino != statSelf.Sys().(*syscall.Stat_t).Ino {
		return fmt.Errorf("inodes don't match, something else must be updating already")
	}

	return nil
}
