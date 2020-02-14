package notes

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBody(t *testing.T) {
	r := body("testing", time.Date(2012, 12, 12, 12, 12, 12, 0, time.UTC))
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		t.Fatalf("Couldn't couldn't copy body: %s", err)
	}

	assert.Equal(t, `{
"title": "testing"
"date": "2012-12-12T12:12:12"
"tags": [ "private", "inbox" ]
}
 * testing
`, buf.String())
}
