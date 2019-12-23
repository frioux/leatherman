package email

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

var update = flag.Bool("update", false, "update golden files")

func TestToJSON(t *testing.T) {
	t.Parallel()

	tests := []string{"basic"}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			sourceMIME, err := os.Open(filepath.Join("testdata", name+".eml"))
			if err != nil {
				t.Fatalf("Couldn't load MIME: %s", err)
			}
			defer sourceMIME.Close()

			buf := &bytes.Buffer{}
			err = toJSON(sourceMIME, buf)
			if err != nil {
				t.Errorf("%s errored: %s", name, err)
				return
			}

			golden := filepath.Join("testdata", name+".json")
			if *update {
				if err := ioutil.WriteFile(golden, buf.Bytes(), 0644); err != nil {
					t.Fatalf("Couldn't update %s: %s", golden, err)
				}
			}
			expected, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatalf("Couldn't load JSON: %s", err)
			}

			testutil.JSONEqual(t, buf.String(), string(expected), name+" matches")
		})
	}
}
