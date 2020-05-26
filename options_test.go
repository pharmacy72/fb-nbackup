package fb_nbackup

import (
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCredential(t *testing.T) {
	assert.Implements(t, (*Argument)(nil), (*Credential)(nil))
}

func TestCredentialToArgument(t *testing.T) {
	f := fuzz.New().NilChance(0)
	var cred Credential
	f.Fuzz(&cred)
	args := cred.ToArgument()
	assert.Len(t, args, 8)
	check := func(name, arg string) {
		if arg != "" {
			assert.Contains(t, args, name)
			assert.Contains(t, args, arg)
		}
	}
	check("-USER", cred.User)
	check("-ROLE", cred.Role)
	check("-PASSWORD", cred.Password)
	check("-FETCH_PASSWORD", cred.PasswordFromFile)
}

func TestOption(t *testing.T) {
	cred := &Credential{}
	decompressCmd := "unzip"
	command := "other cmd"
	options := []Option{
		WithCredential(cred),
		WithWriter(os.Stdout),
		WithDecompressCommand(decompressCmd),
		WithDirect(true),
		WithoutTriggers(true),
		WithCommandPath(command),
	}
	m := &Manager{}
	for _, option := range options {
		option(m)
	}
	assert.Equal(t, cred, m.credential)
	assert.Equal(t, os.Stdout, m.output)
	assert.Equal(t, decompressCmd, m.decompressCommand)
	assert.Equal(t, command, m.command)
	assert.True(t, m.noDBTriggers)
	assert.True(t, m.direct)
}
