package fb_nbackup

import (
	"bytes"
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

func TestManager_Common(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executer := mock_fb_nbackup.NewMockexecuter(ctrl)

	manager := NewManager(
		WithCommandPath("cmd"),
		withExecuter(executer))

	db := "database.fdb"

	commands := []struct {
		name         string
		baseCommand  string
		method       func(context.Context) error
		getArguments func(name string) []interface{}
	}{
		{
			name:        "FIXUP",
			baseCommand: "-FIXUP",
			method: func(ctx context.Context) error {
				return manager.Fixup(ctx, db)
			},
		},
		{
			name:        "UNLOCK",
			baseCommand: "-UNLOCK",
			method: func(ctx context.Context) error {
				return manager.Unlock(ctx, db)
			},
		},
		{
			name:        "BACKUP",
			baseCommand: "-BACKUP",
			method: func(ctx context.Context) error {
				return manager.Backup(ctx, NewLevel(2), db, "some file")
			},
			getArguments: func(name string) []interface{} {
				return []interface{}{name, NewLevel(2).String(), db, "some file"}
			},
		},
		{
			name:        "RESTORE",
			baseCommand: "-RESTORE",
			method: func(ctx context.Context) error {
				return manager.Restore(ctx, db, "file1", "fileA", "fileN")
			},
			getArguments: func(name string) []interface{} {
				return []interface{}{name, db, "file1", "fileA", "fileN"}
			},
		},
	}

	for _, command := range commands {
		t.Run(command.name, func(t *testing.T) {
			getArgs := command.getArguments
			if getArgs == nil {
				getArgs = func(name string) []interface{} {
					return []interface{}{name, db}
				}
			}

			errWant := errors.New("error exec")
			executer.EXPECT().Exec(ctx, "cmd", getArgs(command.baseCommand)...).Return(errWant)
			err := command.method(ctx)
			assert.EqualError(t, err, errWant.Error())

			executer.EXPECT().Exec(ctx, "cmd", getArgs(command.baseCommand)...)
			err = command.method(ctx)
			assert.NoError(t, err)
		})
	}
}

func TestManager_BackupTo(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executer := mock_fb_nbackup.NewMockexecuter(ctrl)

	manager := NewManager(
		WithCommandPath("cmd"),
		withExecuter(executer))

	db := "database.fdb"
	level := NewLevel(2)
	writer := &bytes.Buffer{}

	errWant := errors.New("error exec")
	executer.EXPECT().ExecWithWriter(ctx, "cmd", writer, "-BACKUP", level.String(), db, "stdout").
		Return(errWant)
	err := manager.BackupTo(ctx, level, db, writer)
	assert.EqualError(t, err, errWant.Error())

	executer.EXPECT().ExecWithWriter(ctx, "cmd", writer, "-BACKUP", level.String(), db, "stdout")
	err = manager.BackupTo(ctx, level, db, writer)
	assert.NoError(t, err)
}

func TestManager_Lock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executer := mock_fb_nbackup.NewMockexecuter(ctrl)

	manager := NewManager(
		WithCommandPath("cmd"),
		withExecuter(executer))

	db := "database.fdb"

	errWant := errors.New("error exec")
	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-LOCK", db).Return(errWant)
	_, err := manager.Lock(ctx, db, false)
	assert.EqualError(t, err, errWant.Error())

	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-LOCK", db)
	got, err := manager.Lock(ctx, db, false)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), got)

	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-SIZE", "-LOCK", db).DoAndReturn(
		func(ctx context.Context, commandLine string, w io.Writer, args ...string) error {
			w.Write([]byte(`abc`)) // nolint: errcheck
			return nil
		})
	_, err = manager.Lock(ctx, db, true)
	assert.EqualError(t, err, `strconv.Atoi: parsing "abc": invalid syntax`)
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
			w.Write([]byte(`1.2.3.4`)) // nolint: errcheck
			return nil
		})
	got, err = manager.Version(ctx)
	assert.NoError(t, err)
	assert.Empty(t, got)

	executer.EXPECT().ExecWithWriter(ctx, "cmd", gomock.Any(), "-Z").DoAndReturn(
		func(ctx context.Context, commandLine string, w io.Writer, args ...string) error {
			w.Write([]byte(`nbackup version:V3.0.5.33220`)) // nolint: errcheck
			return nil
		})
	got, err = manager.Version(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "V3.0.5.33220", got)
}
