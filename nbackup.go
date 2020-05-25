package fb_nbackup

import (
	"context"
	"io"
	"strconv"
)

// Backuper is the interface for operations with the nbackup command
type Backuper interface {
	Version(ctx context.Context) (string, error)
	Lock(ctx context.Context, db string, returnSize bool) (int64, error)
	Unlock(ctx context.Context, db string) error
	Fixup(ctx context.Context, db string) error
	Backup(ctx context.Context, level Level, db string, file string) error
	BackupTo(ctx context.Context, level Level, db string, w io.Writer) error
	Restore(ctx context.Context, db string, files ...string) error
}

// Argument is the interface for returning a list of command arguments
type Argument interface {
	ToArgument() []string
}

// The type is determined by the backup level. The valid level starts at zero.
type Level int

// NewLevel returns new level
func NewLevel(i int) Level {
	return Level(i)
}

// ToArgument implements Argument
func (l Level) ToArgument() []string {
	return []string{strconv.Itoa(int(l))}
}

// Int returns level as int
func (l Level) Int() int {
	return int(l)
}

// String returns level as string.
// Implements String interface
func (l Level) String() string {
	return strconv.Itoa(l.Int())
}
