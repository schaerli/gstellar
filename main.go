package main

import (
	"fmt"
	"os"

	"github.com/schaerli/gstellar/initialize"
	"github.com/schaerli/gstellar/snapshot"
	"github.com/schaerli/gstellar/web"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]

		if command == "snapshot" {
			if len(os.Args) > 2 {
				subCommand := os.Args[2]
				if subCommand == "create" {
					snapshot.SnapshotCreatePrepare()
					os.Exit(0)
				}
				if subCommand == "restore" {
					snapshot.SnapshotRestorePrepare()
					os.Exit(0)
				}
				if subCommand == "list" {
					snapshot.SnapshotList()
					os.Exit(0)
				}
				if subCommand == "drop" {
					snapshot.DropSnapshotPrepare()
					os.Exit(0)
				}
			} else {
				fmt.Println("Snapshots Commands:")
				fmt.Println("  create")
				fmt.Println("  restore")
				fmt.Println("  list")
				fmt.Println("  drop")
				os.Exit(0)
			}
			fmt.Println("snapshot arg")
			os.Exit(0)
		}

		if command == "init" {
			initialize.Init()
		}

		if command == "init-only-db" {
			initialize.InitDb()
		}

		if command == "web" {
			web.Start()
		}

	} else {
		fmt.Println("Commands:")
		fmt.Println("  init")
		fmt.Println("  snapshot")
		fmt.Println("  web")
		os.Exit(0)
	}
}

