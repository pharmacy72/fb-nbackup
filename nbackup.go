package fb_nbackup

import (
	"context"
	"io"
)

type Backuper interface {
	Lock(ctx context.Context, db string) error
	Unlock(ctx context.Context, db string) error
	Fixup(ctx context.Context, db string) error
	Backup(ctx context.Context, level int, db string) error
	BackupTo(ctx context.Context, level int, db string, w io.Writer)
	Restore(ctx context.Context, db string, files ...string) error

	Version() (string, error)
}
