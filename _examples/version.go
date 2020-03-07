package main

import (
	nbackup "fb-nbackup"
	"fmt"
	"os"
)

func main() {
	manager := nbackup.NewManager(
		nbackup.WithCommandPath("/usr/bin/docker run --rm jacobalberty/firebird:3.0 /usr/local/firebird/bin/nbackup"),
	//	nbackup.WithCommandPath("ls -l"),
	)
	version, err := manager.Version()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("nbackup version:", version)
}
