package main

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	if _, err := os.Stat("gstellar.yml"); err == nil {
		// path/to/whatever exists

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Enter PG Superuser name: ")
		var superUserName string
		fmt.Scanln(&superUserName)

		fmt.Println("Enter PG Superuser password: ")
		var superUserPass string
		fmt.Scanln(&superUserPass)

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

	}

	fmt.Println("Enter PG Superuser name: ")
	var superUserName string
	fmt.Scanln(&superUserName)

	fmt.Println("Enter PG Superuser password: ")
	var superUserPass string
	fmt.Scanln(&superUserPass)

	fmt.Print(first + " " + second)
	dsn := "host=localhost user=gorm password=gorm dbname=postgres port=5432"
	_db, _err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

}
