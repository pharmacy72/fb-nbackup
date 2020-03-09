package main

import (
	"context"
	nbackup "fb-nbackup"
	"fmt"
	"os"
)

func main() {
	manager := nbackup.NewManager(
		nbackup.WithCommandPath("/usr/bin/docker exec firebird /usr/local/firebird/bin/nbackup"),
	//	nbackup.WithCommandPath("ls -l"),
	)
	version, err := manager.Version()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("nbackup version:", version)

	ctx := context.Background()
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

}
