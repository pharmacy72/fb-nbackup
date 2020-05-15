package fb_nbackup

import (
	"context"
	"io"
	"strconv"
)

type Level int

type Backuper interface {
	Version(ctx context.Context) (string, error)
	Lock(ctx context.Context, db string, returnSize bool) (int64, error)
	Unlock(ctx context.Context, db string) error
	Fixup(ctx context.Context, db string) error
	Backup(ctx context.Context, level Level, db string, file string) error
	BackupTo(ctx context.Context, level Level, db string, w io.Writer) error
	Restore(ctx context.Context, db string, files ...string) error
}

type Argument interface {
	ToArgument() []string
}

func NewLevel(i int) Level {
	return Level(i)
}

func (l Level) ToArgument() []string {
	return []string{strconv.Itoa(int(l))}
}

func (l Level) Int() int {
	return int(l)
}
func (l Level) String() string {
	return strconv.Itoa(l.Int())
}
