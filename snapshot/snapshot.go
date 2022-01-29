package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schaerli/gstellar/initialize"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Snapshot struct {
	Id int
	SnapshottedDb string
	SnapshotName string
	OriginalDb string
	OriginalOwner string
	CreatedAt	time.Time
}

func SnapshotCreatePrepare() {
	db := GetDb()

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

	SnapshotCreate(choosenDb, snapshotName)
}

func SnapshotCreate(choosenDb string, snapshotName string) {
	db := GetDb()
	snapshotDbName := buildSnapshotDbName(choosenDb, snapshotName)
	orignalDbOwner := getDbOwner(db, choosenDb)
	createSnapshotRecord(db, snapshotDbName, snapshotName, choosenDb, orignalDbOwner)
	createSnapshot(db, snapshotDbName, choosenDb, orignalDbOwner)

	output := fmt.Sprintf("Snapshot %s from %s created", snapshotName, choosenDb)
	fmt.Println(output)
}

func SnapshotList() {
	db := GetDb()

	rows, _ := db.Order("id desc").Model(&Snapshot{}).Rows()

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

func SnapshotRestore() {
	db := GetDb()

	var snapshots []Snapshot
	db.Order("id desc").Select("Id", "SnapshotName", "OriginalDb", "CreatedAt").Find(&snapshots)

	var snapshotLabels []string
	for _, s := range snapshots {
		label := fmt.Sprintf("%s | %s | %s | %s", strconv.Itoa(s.Id), s.SnapshotName, s.OriginalDb, s.CreatedAt)
		snapshotLabels = append(snapshotLabels, label)
	}

	choosenSnapshot := ""
	prompt := &survey.Select{
			Message: "Which Snapshot?",
			Options: snapshotLabels,
	}
	survey.AskOne(prompt, &choosenSnapshot, survey.WithValidator(survey.Required))

	r, _ := regexp.Compile(`^\d*`)
	id := r.FindString(choosenSnapshot)

	var snapshot Snapshot
	db.First(&snapshot, id)

	removeDatabase(db, snapshot.OriginalDb)
	createSnapshot(db, snapshot.OriginalDb, snapshot.SnapshottedDb, snapshot.OriginalOwner)

	output := fmt.Sprintf("Snapshot %s on %s restored", snapshot.SnapshotName, snapshot.OriginalDb)
	fmt.Println(output)
}

func GetDb() *gorm.DB {
	dbCredentials := ReadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=gstellar port=%s", dbCredentials.Host, dbCredentials.SuperUserName, dbCredentials.SuperUserPass, dbCredentials.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
    panic("failed to connect database")
  }

	return db
}

func removeDatabase(db *gorm.DB, database  string) {
	queryTemplate := `
		DROP DATABASE IF EXISTS "%s" WITH (FORCE)
	`
	removeQuery := fmt.Sprintf(queryTemplate, database)
	db.Exec(removeQuery)
}

func createSnapshot(db *gorm.DB, snapshotName string, choosenDb string, originalDbOwner string) {
	dropConnectionsTemplate := `
	SELECT pg_terminate_backend(pg_stat_activity.pid)
	FROM pg_stat_activity
	WHERE pg_stat_activity.datname = '%s'
		AND pid <> pg_backend_pid()
	`

	dropConnectionsQuery := fmt.Sprintf(dropConnectionsTemplate, choosenDb)
	db.Exec(dropConnectionsQuery)


	queryTemplate := `
	CREATE DATABASE "%s"
	WITH TEMPLATE "%s"
	OWNER "%s"
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

func ReadConfig() initialize.DbCredentials {
	var dbCredentials initialize.DbCredentials
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