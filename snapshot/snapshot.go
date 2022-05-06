package snapshot

import (
	"fmt"
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
	SizeGb int
}

type SizeGb struct{
	Sizegb int
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
	t.AppendHeader(table.Row{"#", "Snapshot name", "Source db", "Used Giga", "Created at"})

	for rows.Next() {
		var snapshot Snapshot
		db.ScanRows(rows, &snapshot)
		sizeGiga := GetSizeOfDb(db, snapshot.SnapshottedDb)

		t.AppendRow(table.Row{snapshot.Id, snapshot.SnapshotName, snapshot.OriginalDb, sizeGiga, snapshot.CreatedAt})
	}

	t.Render()
}

func SnapshotRestorePrepare() {
	SnapshotRestore(chooseSnapshot())
}

func SnapshotRestore(snapshotId string) string {
	db := GetDb()
	var snapshot Snapshot
	db.First(&snapshot, snapshotId)

	removeDatabase(db, snapshot.OriginalDb)
	createSnapshot(db, snapshot.OriginalDb, snapshot.SnapshottedDb, snapshot.OriginalOwner)
	optimizeRestoredDatabase(snapshot.OriginalDb)

	output := fmt.Sprintf("Snapshot %s on %s restored", snapshot.SnapshotName, snapshot.OriginalDb)
	fmt.Println(output)

	return output
}

func GetDb() *gorm.DB {
	dbCredentials := initialize.ReadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=gstellar port=%s", dbCredentials.Host, dbCredentials.SuperUserName, dbCredentials.SuperUserPass, dbCredentials.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
    panic("failed to connect database")
  }

	return db
}

func removeDatabase(db *gorm.DB, database string) {
	var server_version string
	versionQuery := "SHOW server_version"
	db.Raw(versionQuery).Scan(&server_version)
	floatVar, _ := strconv.ParseFloat(server_version, 32)

	if floatVar <= 13.0 {
		dropConnectionsTemplate := `
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
			AND pid <> pg_backend_pid()
		`

		dropConnectionsQuery := fmt.Sprintf(dropConnectionsTemplate, database)
		db.Exec(dropConnectionsQuery)

		queryTemplate := `
		DROP DATABASE IF EXISTS "%s"
		`
		removeQuery := fmt.Sprintf(queryTemplate, database)
		db.Exec(removeQuery)
	} else {
		queryTemplate := `
		DROP DATABASE IF EXISTS "%s" WITH (FORCE)
		`
		removeQuery := fmt.Sprintf(queryTemplate, database)
		db.Exec(removeQuery)
	}
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

func GetSizeOfDb(db *gorm.DB, dbName string) int {
	getSizeQueryTemplate := `
	SELECT pg_database.datname AS "databasename",
	pg_database_size(pg_database.datname)/1024/1024/1024 AS "sizegb"
	FROM pg_database
	WHERE pg_database.datname='%s'
	`
	dbSizeQuery := fmt.Sprintf(getSizeQueryTemplate, dbName)

	var sizegb SizeGb
	db.Raw(dbSizeQuery).Scan(&sizegb)

	return sizegb.Sizegb
}

func DropSnapshotPrepare() {
	DropSnapshot(chooseSnapshot())
}

func DropSnapshot(id string) string {
	db := GetDb()
	var snapshot Snapshot
	db.First(&snapshot, id)

	removeDatabase(db, snapshot.SnapshottedDb)

	queryTemplate := `
	delete from snapshots where id = '%d'
	`

	dropQuery := fmt.Sprintf(queryTemplate, snapshot.Id)
	db.Exec(dropQuery)
	output := fmt.Sprintf("Snapshot %s dropped database %s", snapshot.SnapshotName, snapshot.SnapshottedDb)
	fmt.Println(output)
	return output
}

func chooseSnapshot() string {
	db := GetDb()

	var snapshots []Snapshot
	db.Order("id desc").Select("Id", "SnapshotName", "OriginalDb", "CreatedAt").Find(&snapshots)

	if len(snapshots) == 0 {
		fmt.Println("No snapshots found!")
		os.Exit(0)
	}

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

	err := survey.AskOne(prompt, &choosenSnapshot, survey.WithValidator(survey.Required))
	switch {
	case err == nil:
		break
	case err.Error() == "interrupt":
		fmt.Println("ctrl-C pressed. Exiting.")
		os.Exit(1)
	case err != nil:
		fmt.Printf("%v. Trying again.\n", err)
	default:
		break
	}

	r, _ := regexp.Compile(`^\d*`)

	return r.FindString(choosenSnapshot)
}

func optimizeRestoredDatabase(dbName string) {
	dbCredentials := initialize.ReadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbCredentials.Host, dbCredentials.SuperUserName, dbCredentials.SuperUserPass, dbName, dbCredentials.Port)
	newDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
    panic("failed to connect database")
  }

	tableTemplateQuery := `
	select t.relname::varchar AS table_name
	FROM pg_class t
	JOIN pg_namespace n ON n.oid = t.relnamespace
	WHERE t.relkind = 'r' and n.nspname::varchar = 'public'
	`
	var dbTables []string
	newDb.Raw(tableTemplateQuery).Scan(&dbTables)

	for _, dbTable := range dbTables {
		analyzeTemplateQuery := "ANALYZE %s"
		analzyQuery := fmt.Sprintf(analyzeTemplateQuery, dbTable)

		newDb.Exec(analzyQuery)
	}
}