package netrc

import (
	"testing"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/login.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}
	heroku, _ := f.Machine("api.heroku.com")
	testutil.Equal(t, heroku.Login, "jeff@heroku.com", "wrong login")
	testutil.Equal(t, heroku.Password, "foo", "wrong password")

	heroku2, _ := f.MachineAndLogin("api.heroku.com", "jeff2@heroku.com")
	testutil.Equal(t, heroku2.Login, "jeff2@heroku.com", "wrong login")
	testutil.Equal(t, heroku2.Password, "bar", "wrong password")
}

func TestSampleMulti(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/sample_multi.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}
	m, _ := f.Machine("m")
	n, _ := f.Machine("n")
	testutil.Equal(t, m.Login, "lm", "wrong login")
	testutil.Equal(t, m.Password, "pm", "wrong password")
	testutil.Equal(t, n.Login, "ln", "wrong login")
	testutil.Equal(t, n.Password, "pn", "wrong password")
}

func TestSampleMultiWithDefault(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/sample_multi_with_default.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	m, _ := f.Machine("m")
	n, _ := f.Machine("n")
	testutil.Equal(t, m.Login, "lm", "wrong login")
	testutil.Equal(t, m.Password, "pm", "wrong password")
	testutil.Equal(t, n.Login, "ln", "wrong login")
	testutil.Equal(t, n.Password, "pn", "wrong password")
}

func TestNewlineless(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/newlineless.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	m, _ := f.Machine("m")
	testutil.Equal(t, m.Login, "l", "wrong login")
	testutil.Equal(t, m.Password, "p", "wrong password")
}

func TestBadDefaultOrder(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/bad_default_order.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	g, _ := f.Machine("mail.google.com")
	r, _ := f.Machine("ray")
	testutil.Equal(t, g.Login, "joe@gmail.com", "wrong login")
	testutil.Equal(t, g.Password, "somethingSecret", "wrong password")
	testutil.Equal(t, r.Login, "demo", "wrong login")
	testutil.Equal(t, r.Password, "mypassword", "wrong password")
}

func TestDefaultOnly(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/default_only.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	d, _ := f.Machine("default")
	testutil.Equal(t, d.Login, "ld", "wrong login")
	testutil.Equal(t, d.Password, "pd", "wrong password")
}

func TestGood(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/good.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	g, _ := f.Machine("mail.google.com")
	testutil.Equal(t, g.Login, "joe@gmail.com", "wrong login")
	testutil.Equal(t, g.Account, "justagmail", "wrong account")
	testutil.Equal(t, g.Password, "somethingSecret", "wrong password")
}

func TestPassword(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/password.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	m, _ := f.Machine("m")
	testutil.Equal(t, m.Password, "p", "wrong password")
}

func TestPermissive(t *testing.T) {
	t.Parallel()

	f, err := Parse("./testdata/permissive.netrc")
	if err != nil {
		t.Errorf("Parse failed: %s", err)
		return
	}

	m, _ := f.Machine("m")
	testutil.Equal(t, m.Login, "l", "wrong login")
	testutil.Equal(t, m.Password, "p", "wrong password")
}
