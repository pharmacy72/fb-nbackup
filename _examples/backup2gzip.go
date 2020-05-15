package main

import (
	"compress/gzip"
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
			User:     "fbuser",
			Password: "023RsdTf4UI123",
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fbackup, err := os.Create("./backup/NBEXAMPLE_L0.fbk.gz")
	if err != nil {
		panic(nil)
	}
	defer fbackup.Close()
	zipWriter := gzip.NewWriter(fbackup)
	defer zipWriter.Close()

	level := nbackup.NewLevel(0)
	dsn := "NBEXAMPLE"
	err = manager.BackupTo(ctx, level, dsn, zipWriter)
	if err != nil {
		fmt.Println("failed to backup", err)
		os.Exit(1)
	}
	fmt.Printf("Backups created. Level: %v\n", level)
}
