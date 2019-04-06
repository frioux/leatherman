package netrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	f, err := Parse("./testdata/login.netrc")
	assert.Nil(t, err)
	heroku := f.Machine("api.heroku.com")
	assert.Equal(t, "jeff@heroku.com", heroku.Login)
	assert.Equal(t, "foo", heroku.Password)

	heroku2 := f.MachineAndLogin("api.heroku.com", "jeff2@heroku.com")
	assert.Equal(t, heroku2.Login, "jeff2@heroku.com")
	assert.Equal(t, heroku2.Password, "bar")
}

func TestSampleMulti(t *testing.T) {
	f, err := Parse("./testdata/sample_multi.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Login, "lm")
	assert.Equal(t, f.Machine("m").Password, "pm")
	assert.Equal(t, f.Machine("n").Login, "ln")
	assert.Equal(t, f.Machine("n").Password, "pn")
}

func TestSampleMultiWithDefault(t *testing.T) {
	f, err := Parse("./testdata/sample_multi_with_default.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Login, "lm")
	assert.Equal(t, f.Machine("m").Password, "pm")
	assert.Equal(t, f.Machine("n").Login, "ln")
	assert.Equal(t, f.Machine("n").Password, "pn")
}

func TestNewlineless(t *testing.T) {
	f, err := Parse("./testdata/newlineless.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Login, "l")
	assert.Equal(t, f.Machine("m").Password, "p")
}

func TestBadDefaultOrder(t *testing.T) {
	f, err := Parse("./testdata/bad_default_order.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("mail.google.com").Login, "joe@gmail.com")
	assert.Equal(t, f.Machine("mail.google.com").Password, "somethingSecret")
	assert.Equal(t, f.Machine("ray").Login, "demo")
	assert.Equal(t, f.Machine("ray").Password, "mypassword")
}

func TestDefaultOnly(t *testing.T) {
	f, err := Parse("./testdata/default_only.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("default").Login, "ld")
	assert.Equal(t, f.Machine("default").Password, "pd")
}

func TestGood(t *testing.T) {
	f, err := Parse("./testdata/good.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("mail.google.com").Login, "joe@gmail.com")
	assert.Equal(t, f.Machine("mail.google.com").Account, "justagmail")
	assert.Equal(t, f.Machine("mail.google.com").Password, "somethingSecret")
}

func TestPassword(t *testing.T) {
	f, err := Parse("./testdata/password.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Password, "p")
}

func TestPermissive(t *testing.T) {
	f, err := Parse("./testdata/permissive.netrc")
	assert.Nil(t, err)
	assert.Equal(t, f.Machine("m").Login, "l")
	assert.Equal(t, f.Machine("m").Password, "p")
}
