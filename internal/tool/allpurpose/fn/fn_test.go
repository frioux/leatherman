package fn

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

type runTest struct {
	name, scriptName string
	in               []string
	out              string
}

func TestRun(t *testing.T) {
	t.Parallel()

	var err error
	dir, err = ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't setup test dir: %s", err)
	}
	defer os.RemoveAll(dir)

	tests := []runTest{
		{
			name: "basic",
			out:  "this is a test\n",
			in:   []string{"basic", `echo "this is a test"`},
		},
		{
			name:       "replace",
			scriptName: "basic",
			out:        "replaced\n",
			in:         []string{"basic", "-f", `echo "replaced"`},
		},
		{
			name: "tokenized",
			out:  "foo bar\n",
			in:   []string{"tokenized", "echo", "foo", "bar"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := []string{"fn"}
			input = append(input, test.in...)
			Run(input, nil)
			name := test.scriptName
			if name == "" {
				name = test.name
			}
			cmd := exec.Command(dir + "/" + name)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("Couldn't run command: %s", err)
			}

			testutil.Equal(t, string(out), test.out, "wrong output")
		})
	}
}
