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
		nbackup.WithCommandPath("/usr/bin/docker exec -w /backup nbackup_fb /usr/local/firebird/bin/nbackup"),
		nbackup.WithCredential(&nbackup.Credential{
			User:     "fbuser",
			Password: "023RsdTf4UI123",
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// levels: 0->1->2
	// output:
	//    Backups created. Level: 0 File: file_l0.fbk
	//    Backups created. Level: 1 File: file_l1.fbk
	//    Backups created. Level: 2 File: file_l2.fbk
	//
	//	  see: ls -la ./backup


	for i := 0; i < 3; i++ {
		level := nbackup.Level(i)
		dsn := "NBEXAMPLE"
		fileBackup := fmt.Sprintf("file_l%s.fbk", level)
		err := manager.Backup(ctx, level, dsn, fileBackup)
		if err != nil {
			fmt.Println("failed to backup", err)
			os.Exit(1)
		}
		fmt.Printf("Backups created. Level: %v File: %s\n", level, fileBackup)
	}
}
