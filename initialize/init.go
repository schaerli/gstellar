package initialize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/AlecAivazis/survey/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbCredentials struct {
	SuperUserName string
	SuperUserPass string
	Host string
	Port string
}

func Init() {
	superUserName := ""
	superUserNameInput := &survey.Input{
			Message: "Enter PG Superuser name:",
	}
	survey.AskOne(superUserNameInput, &superUserName, survey.WithValidator(survey.Required))

	superUserPass := ""
	superUserPassInput := &survey.Password{
			Message: "Enter PG Superuser password:",
	}
	survey.AskOne(superUserPassInput, &superUserPass)

	host := ""
	hostInput := &survey.Input{
			Message: "Enter PG Host:",
			Help: "maybe localhost",
	}
	survey.AskOne(hostInput, &host)

	port := ""
	portInput := &survey.Input{
			Message: "Enter PG Port:",
			Help: "maybe 5432",
	}
	survey.AskOne(portInput, &port)

	jsonObj := DbCredentials{SuperUserName: superUserName, SuperUserPass: superUserPass, Host: host, Port: port}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s", host, superUserName, superUserPass, port)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	var dbName *string
	db.Raw("SELECT datname FROM pg_database WHERE datname = 'gstellar'").Scan(&dbName)

	if dbName == nil {
		fmt.Println("Create gstellar db")
		db.Exec("CREATE DATABASE gstellar")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=gstellar port=%s", host, superUserName, superUserPass, port)
		gstellarDb, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		snapshotTable := `CREATE TABLE IF NOT EXISTS snapshots (
			id             integer not null primary key,
			snapshotted_db varchar(200) not null,
			snapshot_name  varchar(200) not null,
			original_db    varchar(200) not null,
			original_owner varchar(200) not null,
			created_at     timestamp without time zone,
			updated_at     timestamp without time zone
		)`

		gstellarDb.Exec(snapshotTable)

		sequence := `CREATE SEQUENCE IF NOT EXISTS snapshot_sequence
			start 1
			increment 1
		`
		gstellarDb.Exec(sequence)
	}

	file, _ := json.MarshalIndent(jsonObj, "", " ")
	ioutil.WriteFile("gstellar.json", file, 0644)
	fmt.Println("init success")
}
