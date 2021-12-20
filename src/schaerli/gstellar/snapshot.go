package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Snapshot struct {
	Id int
	SnapshottedDb string
	SnapshotName string
	OriginalDb string
	OriginalOwner string
	CreatedAt	time.Time
}

func SnapshotCreate() {
	db := getDb()

	var dbNames []string
	db.Raw("SELECT datname FROM pg_database").Scan(&dbNames)

	choosenDb := ""
	prompt := &survey.Select{
			Message: "Which DB?",
			Options: dbNames,
	}
	survey.AskOne(prompt, &choosenDb, survey.WithValidator(survey.Required))

	snapshotName := ""
	promptInput := &survey.Input{
			Message: "Snapshot name?",
	}
	survey.AskOne(promptInput, &snapshotName)

	snapshotDbName := buildSnapshotDbName(choosenDb, snapshotName)
	orignalDbOwner := getDbOwner(db, choosenDb)
	createSnapshotRecord(db, snapshotDbName, snapshotName, choosenDb, orignalDbOwner)
	createSnapshot(db, snapshotDbName, choosenDb, orignalDbOwner)
}

func SnapshotList() {
	db := getDb()

	rows, _ := db.Model(&Snapshot{}).Rows()
	defer rows.Close()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Snapshot name", "Source db", "Created at"})

	for rows.Next() {
		var snapshot Snapshot
		db.ScanRows(rows, &snapshot)

		t.AppendRow(table.Row{snapshot.Id, snapshot.SnapshotName, snapshot.OriginalDb, snapshot.CreatedAt})
	}

	t.Render()
}

func getDb() *gorm.DB {
	dbCredentials := ReadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=gstellar port=%s", dbCredentials.Host, dbCredentials.SuperUserName, dbCredentials.SuperUserPass, dbCredentials.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
    panic("failed to connect database")
  }

	return db
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