package main_test

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"
)

func TestNoPrintln(t *testing.T) {
	cmd := exec.Command("git", "grep", "--quiet", "--perl-regexp", `\bpr[i]nt(ln)?\(`, "*.go")
	err := cmd.Run()
	var e *exec.ExitError

	if errors.As(err, &e) && e.ExitCode() != 0 {
		return
	}

	if err != nil {
		t.Errorf("unexpected error from git grep: %s\n", err)
	} else {
		b := &bytes.Buffer{}
		t.Errorf("git grep found the forbidden println:")
		cmd := exec.Command("git", "grep", "--perl-regexp", `\bpr[i]nt(ln)?\(`, "*.go")
		cmd.Stdout = b
		cmd.Stderr = b
		cmd.Run()
		t.Error(b.String())
	}
}
