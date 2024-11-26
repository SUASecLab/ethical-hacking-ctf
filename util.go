package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/SUASecLab/workadventure_admin_extensions/extensions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func isAuthenticated(w http.ResponseWriter, token string) bool {
	validationResult, err := extensions.GetValidationResult("http://" + sidecarUrl + "/validate?token=" + token)

	if err != nil || !validationResult.Valid {
		fmt.Fprintln(w, "Invalid authentication token")
		return false
	}
	return true
}

func getUuidFromToken(token string) (string, bool) {
	userInfo, err := extensions.GetUserInfo("http://" + sidecarUrl + "/userinfo?token=" + token)
	if err != nil {
		log.Println("Can not receive UUID: ", err)
		return "", false
	}
	return userInfo.UUID, true
}

func getUserByUuid(usersCollection *mongo.Collection, uuid string) (User, error) {
	var user User
	err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}}).Decode(&user)
	return user, err
}

// this should be effectively the last line of code invoked in the calling function
func runTemplate(w http.ResponseWriter, data IndexData) {
	// see https://dev.to/moniquelive/passing-multiple-arguments-to-golang-templates-16h8
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"arr": func(els ...any) []any { return els },
	}).ParseFiles("templates/successMessage.html",
		"templates/errorMessage.html", "templates/progressbar.html", "templates/index.html", "templates/base.html"))
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Could not parse template: ", err)
		fmt.Fprint(w, "Can not parse template to show data")
	}
}
