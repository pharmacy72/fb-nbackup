package main

import (
	"context"
	nbackup "fb-nbackup"
	"fmt"
	"time"
)

func main() {
	manager := nbackup.NewManager(
		nbackup.WithCommandPath("/usr/bin/docker exec nbackup_fb /usr/local/firebird/bin/nbackup"), //uncomment to run in docker
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	version, err := manager.Version(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("version:", version) // version: V3.0.5.33220
}
