package datefmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateFormat(t *testing.T) {
	assert.Equal(t, "2006-01-02", TranslateFormat("%F"))
	assert.Equal(t, "2006-01-02T15:04:05", TranslateFormat("%FT%T"))
}
