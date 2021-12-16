package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() {
	fmt.Println("Enter PG Superuser name: ")
	var superUserName string
	fmt.Scanln(&superUserName)

	fmt.Println("Enter PG Superuser password: ")
	var superUserPass string
	fmt.Scanln(&superUserPass)

	yml := DbCredentials{SuperUserName: superUserName, SuperUserPass: superUserPass}

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=postgres port=5432", superUserName, superUserPass)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	var dbName *string
	db.Raw("SELECT datname FROM pg_database WHERE datname = 'gstellar'").Scan(&dbName)

	if dbName == nil {
		fmt.Println("Create gstellar db")
		db.Exec("CREATE DATABASE gstellar")

		dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=gstellar port=5432", superUserName, superUserPass)
		gstellarDb, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		snapshotTable := `CREATE TABLE IF NOT EXISTS snapshots (
			id          integer not null primary key,
			name        varchar(200) not null,
			original_db varchar(200) not null,
			created_at timestamp without time zone,
			updated_at timestamp without time zone
		)`

		gstellarDb.Exec(snapshotTable)
	}

	file, _ := json.MarshalIndent(yml, "", " ")
	ioutil.WriteFile("gstellar.json", file, 0644)
}