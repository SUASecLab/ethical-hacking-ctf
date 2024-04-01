package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// this should be effectively the last line of code invoked in the index function
func runTemplate(w http.ResponseWriter, data IndexData) {
	//TODO: add progress bar
	// see https://dev.to/moniquelive/passing-multiple-arguments-to-golang-templates-16h8
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"arr": func(els ...any) []any { return els },
	}).ParseFiles("templates/successMessage.html",
		"templates/errorMessage.html", "templates/progressbar.html", "templates/index.html", "templates/base.html"))
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Could not parse template: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	var token string

	switch r.Method {
	case "GET":
		// user requested to show current flags

		// authenticate user
		token = r.URL.Query().Get("token")
		if !isAuthenticated(w, token) {
			fmt.Fprint(w, "Can not authenticate user")
			return
		}

		// show collected flags
		runTemplate(w, createReturnDataStructure(token, "", ""))
	case "POST":
		// user requested to add a new flag

		// parse form
		if err := r.ParseForm(); err != nil {
			log.Println("Could not parse form: ", err)
			return
		}

		// authenticate user
		token = r.FormValue("token")
		if !isAuthenticated(w, token) {
			fmt.Fprint(w, "Can not authenticate user")
			return
		}

		// get input from parsed form
		flagInput := r.FormValue("flagInput")
		flagTypeSelect := r.FormValue("flagTypeSelect")
		flagIdentifier := flagTypeSelect + "-" + flagInput

		// add flag to user account
		addFlagToUser(w, token, flagIdentifier, flagInput)
	}
}
