package netrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	f, err := Parse("./testdata/login.netrc")
	assert.Nil(t, err)
	heroku, _ := f.Machine("api.heroku.com")
	assert.Equal(t, "jeff@heroku.com", heroku.Login)
	assert.Equal(t, "foo", heroku.Password)

	heroku2, _ := f.MachineAndLogin("api.heroku.com", "jeff2@heroku.com")
	assert.Equal(t, heroku2.Login, "jeff2@heroku.com")
	assert.Equal(t, heroku2.Password, "bar")
}

func TestSampleMulti(t *testing.T) {
	f, err := Parse("./testdata/sample_multi.netrc")
	assert.Nil(t, err)
	m, _ := f.Machine("m")
	n, _ := f.Machine("n")
	assert.Equal(t, m.Login, "lm")
	assert.Equal(t, m.Password, "pm")
	assert.Equal(t, n.Login, "ln")
	assert.Equal(t, n.Password, "pn")
}

func TestSampleMultiWithDefault(t *testing.T) {
	f, err := Parse("./testdata/sample_multi_with_default.netrc")
	assert.Nil(t, err)

	m, _ := f.Machine("m")
	n, _ := f.Machine("n")
	assert.Equal(t, m.Login, "lm")
	assert.Equal(t, m.Password, "pm")
	assert.Equal(t, n.Login, "ln")
	assert.Equal(t, n.Password, "pn")
}

func TestNewlineless(t *testing.T) {
	f, err := Parse("./testdata/newlineless.netrc")
	assert.Nil(t, err)
	m, _ := f.Machine("m")
	assert.Equal(t, m.Login, "l")
	assert.Equal(t, m.Password, "p")
}

func TestBadDefaultOrder(t *testing.T) {
	f, err := Parse("./testdata/bad_default_order.netrc")
	assert.Nil(t, err)
	g, _ := f.Machine("mail.google.com")
	r, _ := f.Machine("ray")
	assert.Equal(t, g.Login, "joe@gmail.com")
	assert.Equal(t, g.Password, "somethingSecret")
	assert.Equal(t, r.Login, "demo")
	assert.Equal(t, r.Password, "mypassword")
}

func TestDefaultOnly(t *testing.T) {
	f, err := Parse("./testdata/default_only.netrc")
	assert.Nil(t, err)
	d, _ := f.Machine("default")
	assert.Equal(t, d.Login, "ld")
	assert.Equal(t, d.Password, "pd")
}

func TestGood(t *testing.T) {
	f, err := Parse("./testdata/good.netrc")
	assert.Nil(t, err)
	g, _ := f.Machine("mail.google.com")
	assert.Equal(t, g.Login, "joe@gmail.com")
	assert.Equal(t, g.Account, "justagmail")
	assert.Equal(t, g.Password, "somethingSecret")
}

func TestPassword(t *testing.T) {
	f, err := Parse("./testdata/password.netrc")
	assert.Nil(t, err)
	m, _ := f.Machine("m")
	assert.Equal(t, m.Password, "p")
}

func TestPermissive(t *testing.T) {
	f, err := Parse("./testdata/permissive.netrc")
	assert.Nil(t, err)
	m, _ := f.Machine("m")
	assert.Equal(t, m.Login, "l")
	assert.Equal(t, m.Password, "p")
}
