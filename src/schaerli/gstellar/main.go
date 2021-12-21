package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]

		if command == "snapshot" {
			if len(os.Args) > 2 {
				subCommand := os.Args[2]
				if subCommand == "create" {
					SnapshotCreate()
					os.Exit(0)
				}
				if subCommand == "restore" {
					SnapshotRestore()
					os.Exit(0)
				}
				if subCommand == "list" {
					SnapshotList()
					os.Exit(0)
				}
			} else {
				fmt.Println("Snapshots Commands:")
				fmt.Println("  create")
				fmt.Println("  list")
				os.Exit(0)
			}
			fmt.Println("snapshot arg")
			os.Exit(0)
		}

		if command == "init" {
			Init()
		}

	} else {
		fmt.Println("Commands:")
		fmt.Println("  init")
		fmt.Println("  snapshot")
		os.Exit(0)
	}
}

func ReadConfig() DbCredentials {
	var dbCredentials DbCredentials
	jsonFileName := "gstellar.json"

	if _, err := os.Stat(jsonFileName); err == nil {
		jsonFile, _ := os.Open(jsonFileName)
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &dbCredentials)

		return dbCredentials
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No gstellar.json found here - run 'gstellar init' first")
	}

	return dbCredentials
}