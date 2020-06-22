package fb_nbackup

import (
	"context"
	"errors"
	mock_fb_nbackup "fb-nbackup/mock"
	"github.com/golang/mock/gomock"
	"io"
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

func TestManager_Version(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executer := mock_fb_nbackup.NewMockexecuter(ctrl)

	manager := NewManager(
		WithCommandPath("cmd"),
		withExecuter(executer))

	errWant := errors.New("error exec")
	executer.EXPECT().ExecWithWriter(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errWant)
	_, err := manager.Version(ctx)
	assert.EqualError(t, err, errWant.Error())

	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-Z")
	got, err := manager.Version(ctx)
	assert.NoError(t, err)
	assert.Empty(t, got)

	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-Z").DoAndReturn(
		func(ctx context.Context, commandLine string, w io.Writer, args ...string) error {
			w.Write([]byte(`1.2.3.4`))
			return nil
		})
	got, err = manager.Version(ctx)
	assert.NoError(t, err)
	assert.Empty(t, got)

	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-Z").DoAndReturn(
		func(ctx context.Context, commandLine string, w io.Writer, args ...string) error {
			w.Write([]byte(`nbackup version:V3.0.5.33220`))
			return nil
		})
	got, err = manager.Version(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "V3.0.5.33220", got)
}
