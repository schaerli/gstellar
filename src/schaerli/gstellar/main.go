package main

import (
	"errors"
	"fmt"
	"os"
)

type DbCredentials struct {
	SuperUserName string
	SuperUserPass string
}

func main() {

	if _, err := os.Stat("gstellar.json"); err == nil {
		if len(os.Args) > 1 {
			command := os.Args[1]

			if command == "snapshot" {
				fmt.Println("snapshot arg")
				os.Exit(0)
			}

		} else {
			fmt.Println("Commands:")
			fmt.Println("  snapshot")
			os.Exit(0)
		}

	} else if errors.Is(err, os.ErrNotExist) {
		Init()
	} else {
		fmt.Println("else here")
	}
}
