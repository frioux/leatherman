package genpass

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	buf := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}

	err := run(buf, bufErr, "password", 0)
	assert.NoError(t, err)
	assert.Regexp(t, `^\$2a\$10\$`, buf.String())
	assert.Regexp(t, `elapsed`, bufErr.String())
}
