package fb_nbackup

import (
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestLevel(t *testing.T) {
	assert.Implements(t, (*Argument)(nil), (*Level)(nil))

	f := fuzz.New().NilChance(0)
	var i int
	f.Fuzz(&i)
	level := NewLevel(i)
	str := strconv.Itoa(i)
	assert.Equal(t, i, level.Int())
	assert.Equal(t, str, level.String())
	assert.Contains(t, level.ToArgument(), str)
	assert.Len(t, level.ToArgument(), 1)
}
