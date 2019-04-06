package netrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	f, err := Parse("./testdata/login.netrc")
	assert.Nil(t, err)
	heroku := f.Machine("api.heroku.com")
	assert.Equal(t, "jeff@heroku.com", heroku.Get("login"))
	assert.Equal(t, "foo", heroku.Get("password"))

	heroku2 := f.MachineAndLogin("api.heroku.com", "jeff2@heroku.com")
	assert.Equal(t, heroku2.Get("login"), "jeff2@heroku.com")
	assert.Equal(t, heroku2.Get("password"), "bar")
}

func TestSampleMulti(t *testing.T) {
	f, err := Parse("./testdata/sample_multi.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Get("login"), "lm")
	assert.Equal(t, f.Machine("m").Get("password"), "pm")
	assert.Equal(t, f.Machine("n").Get("login"), "ln")
	assert.Equal(t, f.Machine("n").Get("password"), "pn")
}

func TestSampleMultiWithDefault(t *testing.T) {
	f, err := Parse("./testdata/sample_multi_with_default.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Get("login"), "lm")
	assert.Equal(t, f.Machine("m").Get("password"), "pm")
	assert.Equal(t, f.Machine("n").Get("login"), "ln")
	assert.Equal(t, f.Machine("n").Get("password"), "pn")
}

func TestNewlineless(t *testing.T) {
	f, err := Parse("./testdata/newlineless.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Get("login"), "l")
	assert.Equal(t, f.Machine("m").Get("password"), "p")
}

func TestBadDefaultOrder(t *testing.T) {
	f, err := Parse("./testdata/bad_default_order.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("mail.google.com").Get("login"), "joe@gmail.com")
	assert.Equal(t, f.Machine("mail.google.com").Get("password"), "somethingSecret")
	assert.Equal(t, f.Machine("ray").Get("login"), "demo")
	assert.Equal(t, f.Machine("ray").Get("password"), "mypassword")
}

func TestDefaultOnly(t *testing.T) {
	f, err := Parse("./testdata/default_only.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("default").Get("login"), "ld")
	assert.Equal(t, f.Machine("default").Get("password"), "pd")
}

func TestGood(t *testing.T) {
	f, err := Parse("./testdata/good.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("mail.google.com").Get("login"), "joe@gmail.com")
	assert.Equal(t, f.Machine("mail.google.com").Get("account"), "justagmail")
	assert.Equal(t, f.Machine("mail.google.com").Get("password"), "somethingSecret")
}

func TestPassword(t *testing.T) {
	f, err := Parse("./testdata/password.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Get("password"), "p")
}

func TestPermissive(t *testing.T) {
	f, err := Parse("./testdata/permissive.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Get("login"), "l")
	assert.Equal(t, f.Machine("m").Get("password"), "p")
}
