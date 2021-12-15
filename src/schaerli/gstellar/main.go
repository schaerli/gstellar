package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type DbCredentials struct {
	SuperUserName string
	SuperUserPass string
}

func main() {

	if _, err := os.Stat("gstellar.json"); err == nil {
		fmt.Println("yaml here")

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Enter PG Superuser name: ")
		var superUserName string
		fmt.Scanln(&superUserName)

		fmt.Println("Enter PG Superuser password: ")
		var superUserPass string
		fmt.Scanln(&superUserPass)

		yml := DbCredentials{SuperUserName: superUserName, SuperUserPass: superUserPass}

		file, _ := json.MarshalIndent(yml, "", " ")
		ioutil.WriteFile("gstellar.json", file, 0644)

	} else {
		fmt.Println("else here")
	}

	// dsn := "host=localhost user=gorm password=gorm dbname=postgres port=5432"
	// _db, _err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

}
