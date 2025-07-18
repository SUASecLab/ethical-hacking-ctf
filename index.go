package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var token string

	switch r.Method {
	case "GET":
		// user requested to show current flags

		// authenticate user
		token = r.URL.Query().Get("token")
		token = html.EscapeString(token)
		if !isAuthenticated(w, token) {
			fmt.Fprint(w, "Can not authenticate user")
			return
		}

		// show collected flags
		runTemplate(w, createReturnDataStructure(token, "", ""))
	case "POST":
		// user requested to add a new flag

		// authenticate user
		token = r.FormValue("token")
		token = html.EscapeString(token)
		if !isAuthenticated(w, token) {
			fmt.Fprint(w, "Can not authenticate user")
			return
		}

		// parse form
		if err := r.ParseForm(); err != nil {
			log.Println("Could not parse form: ", err)
			return
		}

		// get input from parsed form
		flagInput := r.FormValue("flagInput")
		flagTypeSelect := r.FormValue("flagTypeSelect")

		// escape input
		flagInput = html.EscapeString(flagInput)
		flagTypeSelect = html.EscapeString(flagTypeSelect)

		// get flag identifier
		flagIdentifier := flagTypeSelect + "-" + flagInput

		// add flag to user account
		addFlagToUser(w, token, flagIdentifier, flagInput)
	}
}
