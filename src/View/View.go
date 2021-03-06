// Author: Erich Ray
// Title: View
// Version 1
// Created on 1/15/2015
// Manages the views for TimeServer
package view

//import packages
import (
	"fmt"
	"Utility"
	"html/template"
	"net/http"
	"path/filepath"
)

//global variable
var (
	webPrint    = fmt.Fprintln
	getTemplate = func(tmpl string) string {
		tplPath := filepath.Join("templates", tmpl)
		utility.WriteTrace(tplPath)
		return tplPath
	}
	templates *template.Template
)

type TimeData struct {
	Name    string
	Message string
	Time    string
	Port	string
}

func CreateSite() {
	templates = template.Must(template.ParseFiles(getTemplate("about.html"),
		getTemplate("head.html"), getTemplate("index.html"),
		getTemplate("login.html"), getTemplate("logout.html"),
		getTemplate("menu.html"), getTemplate("time.html"),
		getTemplate("logo.html")))
}

func ShowPage(response http.ResponseWriter, site string, data TimeData) error {
	//show correct template
	err := templates.ExecuteTemplate(response, site+".html", data)
	return err
}

func ShowError(response http.ResponseWriter, err error) {
	//debug spew
	utility.WriteTrace("Entering errorHandler")

	//write using println to web response in case problem is
	//with templates
	webPrint(response, "<html>")
	webPrint(response, "<head>")
	webPrint(response, "<style> p {font-size: xx-large}")
	webPrint(response, "</style>")
	webPrint(response, "</head>")
	webPrint(response, "<body>")
	webPrint(response, "<p>"+err.Error())
	webPrint(response, "</p>")
	webPrint(response, "</body>")
	webPrint(response, "</html>")
}
