package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func SnapshotCreate() {
	dbCredentials := ReadConfig()
	fmt.Println(dbCredentials)

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=gstellar port=5432", dbCredentials.SuperUserName, dbCredentials.SuperUserPass)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	var dbName []string
	db.Raw("SELECT datname FROM pg_database").Scan(&dbName)

	choosenDb := ""
	prompt := &survey.Select{
			Message: "Which DB?",
			Options: dbName,
	}
	survey.AskOne(prompt, &choosenDb, survey.WithValidator(survey.Required))

	snapshotName := ""
	promptInput := &survey.Input{
			Message: "Snapshot name?",
	}
	survey.AskOne(promptInput, &snapshotName)

	fmt.Println(choosenDb)
	fmt.Println(snapshotName)

	snapshotDbName := buildSnapshotDbName(choosenDb, snapshotName)
	fmt.Println(snapshotDbName)
	orignalDbOwner := getDbOwner(db, choosenDb)
	createSnapshotRecord(db, snapshotDbName, snapshotName, choosenDb, orignalDbOwner)

	createSnapshot(db, snapshotDbName, choosenDb, orignalDbOwner)
}

func createSnapshot(db *gorm.DB, snapshotName string, choosenDb string, originalDbOwner string) {
	queryTemplate := `
	CREATE DATABASE %s
	WITH TEMPLATE %s
	OWNER %s
	`

	ownerQuery := fmt.Sprintf(queryTemplate, snapshotName, choosenDb, originalDbOwner)
	db.Exec(ownerQuery)
}

func buildSnapshotDbName(choosenDb string, snapshotName string) string {
	rand.Seed(time.Now().UnixNano())
	v := rand.Perm(8)

	result := ""
	for _, s := range v {
		result += strconv.FormatInt(int64(s), 10)
	}

	return fmt.Sprintf("gstellar_%s_%s_%s", choosenDb, snapshotName, result)
}

func getDbOwner(db *gorm.DB, originalDb string) string {
	dbOwner := `
	SELECT d.datname as "Name",
	pg_catalog.pg_get_userbyid(d.datdba) as "ownerName"
	FROM pg_catalog.pg_database d
	WHERE d.datname = '%s'
	ORDER BY 1
	`

	ownerQuery := fmt.Sprintf(dbOwner, originalDb)
	var ownerName []string
	db.Raw(ownerQuery).Scan(&ownerName)

	return ownerName[0]
}

func createSnapshotRecord(db *gorm.DB, snapshotDbName string,
	snapshotName string, originalDb string, orignalDbOwner string,
	) {
	queryTemplate := `
	INSERT INTO snapshots
  (id, snapshotted_db, snapshot_name, original_db, original_owner, created_at, updated_at)
	VALUES
  (nextval('snapshot_sequence'), '%s', '%s', '%s', '%s', now(), now())
	`

	insertQuery := fmt.Sprintf(queryTemplate, snapshotDbName, snapshotName, originalDb, orignalDbOwner)
	db.Exec(insertQuery)
}