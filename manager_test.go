package fb_nbackup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	counter := 0
	optionManager := func(m *Manager) {
		counter++
	}
	manager := NewManager(optionManager, optionManager)
	assert.NotNil(t, manager)
	assert.Equal(t, counter, 2)
	assert.Equal(t, manager, manager.executer)
}

type testToArgument struct {
	Slice []string
}

func (t *testToArgument) ToArgument() []string {
	return t.Slice[:]
}

func TestParseArguments(t *testing.T) {
	cases := []struct {
		Arg  interface{}
		Want []string
		Err  string
	}{
		{
			Arg:  "str",
			Want: []string{"str"},
		},
		{
			Arg:  []string{"s1", "s2", "s3"},
			Want: []string{"s1", "s2", "s3"},
		},
		{
			Arg:  &testToArgument{[]string{"s4", "s5", "s6"}},
			Want: []string{"s4", "s5", "s6"},
		},
		{
			Arg: 123,
			Err: "unknown argument type: int",
		},
	}
	for _, test := range cases {
		got, err := parseArguments(test.Arg)
		if test.Err == "" {
			assert.EqualValues(t, got, test.Want)
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.Err)
		}
	}
}
