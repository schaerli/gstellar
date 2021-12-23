package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Names struct{
	PageTitle string
}

func Start() {
	fmt.Println("Listening on :5678")
	http.HandleFunc("/", handler)
	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	log.Fatal(http.ListenAndServe(":5678", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("main.tmpl").ParseFiles("./web/templates/main.tmpl"))

	data := Names{
		PageTitle: "Holy Moly",
	}

	tmpl.Execute(w, data)
}