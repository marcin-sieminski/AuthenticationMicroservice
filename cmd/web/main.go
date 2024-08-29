package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	cfg = config{api: "http://localhost:81"}
	app = &application{
		config: cfg,
	}
)

type config struct {
	api string
}

type application struct {
	config config
}

type templateData struct {
	API string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "home.page.gohtml")
	})

	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", t))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	td := &templateData{API: app.config.api}

	if err := tmpl.Execute(w, td); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
