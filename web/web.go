package web

import (
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

func Start() {
	fmt.Println("Listening on :5678")

	http.HandleFunc("/", handler)
	http.HandleFunc("/snapshots/new", snapshotNew)
	http.HandleFunc("/snapshots/create", snapshotCreate)

	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	log.Fatal(http.ListenAndServe(":5678", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("main.gohtml").ParseFiles("./web/templates/main.gohtml"))

	db := snapshot.GetDb()
	snapshotRows, _ := db.Model(&snapshot.Snapshot{}).Rows()
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



	http.Redirect(w, r, "/", http.StatusSeeOther)
}