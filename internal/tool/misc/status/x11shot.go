package status

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

type x11shot struct {
	pic []byte
}

func (v *x11shot) load() error {
	f, err := ioutil.TempFile("", "*.png")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	cmd := exec.Command("import", "-window", "root", f.Name())
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	if err := cmd.Run(); err != nil {
		return err
	}

	v.pic, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return nil
}

func (v *x11shot) render(rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "image/png")
	io.Copy(rw, bytes.NewBuffer(v.pic))
}
