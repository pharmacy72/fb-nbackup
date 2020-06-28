package main

import (
	"context"
	nbackup "fb-nbackup"
	"fmt"
	"os"
	"time"
)

func main() {
	manager := nbackup.NewManager(
		nbackup.WithCommandPath("/usr/bin/docker exec nbackup_fb /usr/local/firebird/bin/nbackup"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dsn := "NBEXAMPLE"

	if err := manager.Unlock(ctx, dsn); err != nil {
		fmt.Println("warn:", err)
	}

	size, err := manager.Lock(ctx, dsn, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("database size: %d\n", size)

	if err := manager.Unlock(ctx, dsn); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := manager.Fixup(ctx, dsn); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
