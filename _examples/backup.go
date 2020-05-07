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
		nbackup.WithCommandPath("/usr/bin/docker exec -w /backup firebird /usr/local/firebird/bin/nbackup"),
		nbackup.WithCredential(&nbackup.Credential{
			User:             "fbuser",
			Password:         "023RsdTf4UI123",
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// levels: 0->1->2
	// output:
	//    Backups created. Level: 0
	//    Backups created. Level: 1
	//    Backups created. Level: 2

	for i := 0; i < 3; i++ {
		level := nbackup.Level(i)
		dsn := "NBEXAMPLE"
		fileBackup := ""
		err := manager.Backup(ctx, level, dsn, fileBackup)
		if err != nil {
			fmt.Println("failed to backup", err)
			os.Exit(1)
		}
		fmt.Printf("Backups created. Level: %v\n", level)
	}
}
