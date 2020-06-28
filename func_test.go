//+build functest

package fb_nbackup

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFunctional(t *testing.T) {
	timeout := time.Second * 5; //time.Minute *  5;
	//TODO: WithWriter
	manager := NewManager(
		WithCommandPath("/usr/bin/docker exec -w /backup nbackup_fb /usr/local/firebird/bin/nbackup"),
		WithCredential(&Credential{
			User:     "fbuser",
			Password: "023RsdTf4UI123",
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan bool)
	go func() {
		dsn := "NBEXAMPLE"

		defer close(done)
		version, err := manager.Version(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, version)

		size, err := manager.Lock(ctx, dsn, true)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotZero(t, size)

		if err := manager.Unlock(ctx, dsn); err != nil {
			t.Fatal(err)
		}

		_, err = manager.Lock(ctx, dsn, false)
		if err != nil {
			t.Fatal(err)
		}
		if err := manager.Fixup(ctx, dsn); err != nil {
			t.Fatal(err)
		}

		var files []string

		for i := 0; i < 3; i++ {
			level := Level(i)
			fileBackup := fmt.Sprintf("test_file_l%s.fbk", level)
			files = append(files, fileBackup)
			err := manager.Backup(ctx, level, dsn, fileBackup)
			if err != nil {
				t.Fatal(err)
			}
		}

		dsnDest := "NBEXAMPLE2"
		err = manager.Restore(ctx, dsnDest, files...)
		if err != nil {
			t.Fatal(err)
		}

		buf := &bytes.Buffer{}
		level := NewLevel(0)
		err = manager.BackupTo(ctx, level, dsn, buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotZero(t, buf.Len())
	}()

	select {
	case <-done:
	case <-ctx.Done():
		if ctx.Err() != nil {
			t.Error(ctx.Err())
		}
	}

}
