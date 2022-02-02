package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/schaerli/gstellar/snapshot"
)

type snapshotNewData struct {
	DbNames []string
}

type indexData struct{
	Snapshots []snapshot.Snapshot
}

type restoreJson struct{
	Message string
}

func Start() {
	fmt.Println("Listening on :5678")

	http.HandleFunc("/", handler)
	http.HandleFunc("/snapshots/new", snapshotNew)
	http.HandleFunc("/snapshots/create", snapshotCreate)
	http.HandleFunc("/snapshots/restore", snapshotResore)

	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	log.Fatal(http.ListenAndServe(":5678", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("main.gohtml").ParseFiles("./web/templates/main.gohtml"))

	db := snapshot.GetDb()
	snapshotRows, _ := db.Order("id desc").Model(&snapshot.Snapshot{}).Rows()
	defer snapshotRows.Close()

	snapshots := make([]snapshot.Snapshot, 0)

	for snapshotRows.Next() {
		var snapshot snapshot.Snapshot
		db.ScanRows(snapshotRows, &snapshot)

		snapshots = append(snapshots, snapshot)
	}

	data := &indexData{
		Snapshots: snapshots,
	}

	tmpl.Execute(w, data)
}

func snapshotNew(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("snapshotNew.gohtml").ParseFiles("./web/templates/snapshotNew.gohtml"))

	db := snapshot.GetDb()

	var dbNames []string
	db.Raw("SELECT datname FROM pg_database").Scan(&dbNames)

	data := &snapshotNewData{
		DbNames: dbNames,
	}

	tmpl.Execute(w, data)
}

func snapshotCreate(w http.ResponseWriter, r *http.Request) {
	choosenDB := r.FormValue("selectDB")
	snapshotName := r.FormValue("snapshotName")

	snapshot.SnapshotCreate(choosenDB, snapshotName)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func snapshotResore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	snapshotId, ok := r.URL.Query()["snapshot_id"]

	if !ok || len(snapshotId[0]) < 1 {
			log.Println("Url Param 'snapshot_id' is missing")
			return
	}

	msg := snapshot.SnapshotRestore(snapshotId[0])
	response := restoreJson{Message: msg}
	json.NewEncoder(w).Encode(response)
}