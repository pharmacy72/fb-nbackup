package main

import (
	"context"
	nbackup "fb-nbackup"
	"fmt"
	"os"
	"time"
)

func main() {
	//use: go run _example/backup2stdout.go > level0.fbk
	manager := nbackup.NewManager(
		nbackup.WithCommandPath("/usr/bin/docker exec -w /backup firebird /usr/local/firebird/bin/nbackup"),
		nbackup.WithCredential(&nbackup.Credential{
			User:     "fbuser",
			Password: "023RsdTf4UI123",
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	level := nbackup.NewLevel(0)
	dsn := "NBEXAMPLE"
	err := manager.BackupTo(ctx, level, dsn, os.Stdout)
	if err != nil {
		fmt.Println("failed to backup", err)
		os.Exit(1)
	}
	fmt.Printf("Backups created. Level: %v\n", level)
}
