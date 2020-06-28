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

	dsn := "NBEXAMPLE2"
	files := []string{"file_l0.fbk", "file_l1.fbk", "file_l2.fbk"}
	err := manager.Restore(ctx, dsn, files...)
	if err != nil {
		fmt.Println("failed to backup", err)
		os.Exit(1)
	}
	fmt.Printf("Restored %s.\n", dsn)

	//see:
	// docker exec -i firebird /usr/local/firebird/bin/isql -user fbuser -password 023RsdTf4UI123 << EOF
	// connect NBEXAMPLE;
	// select * from rdb$database;
	// EOF
}
