// Author: Erich Ray
// Title: TimeServer
// Version 3
// Created on 1/15/2015
// Creates a web server for the url http://localhost:8080/
package main

//import packages
import (
	"flag"
	"fmt"
	"github.com/TimeServer/Auth"
	"github.com/TimeServer/Utility"
	"github.com/TimeServer/View"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const version = "Version: 3"

//global variable
var (
	verbose *bool
	port    *int
	debug   *bool
	log     *string
)

//main function for Time Server
func main() {
	var err error

	//parse flags
	verbose = flag.Bool("v", false, "a bool")
	port = flag.Int("port", 8080, "port for webserver")
	log = flag.String("log", "", "load location for seelog")
	debug = flag.Bool("debug", false, "turn optional for debug spew")
	utility.Debug(*debug)
	flag.Parse()

	//build template tree
	view.CreateSite()

	//check if log was specified
	if *log == "" {
		utility.WriteInfo("No log config xml loaded, outputing to console only.")
	} else {

		//load config file and check for errors
		err = utility.SetLogConfig(*log)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	//defer a flush for on quit
	defer utility.FlushLog()

	//print debug
	utility.WriteTrace("Starting Main")

	if *verbose == true {
		fmt.Println(version)
		os.Exit(0)
	}

	//write version number to log
	utility.WriteInfo(version)

	//setup handler for web page requests
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/time/", timeHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/about/", aboutHandler)
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles/"))))

	//set listen and serve for port, checking for error
	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	//if error was returned, send error to standard output
	if err != nil {
		utility.WriteCritical(err.Error())
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

//Function for handling the calls to the about page
func aboutHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("Entering aboutHandler")

	//create data structure
	data := view.TimeData{
		"", "", "",
	}

	//execute template
	err := view.ShowPage(response, "about", data)

	//process error
	if err != nil {
		errorHandler(response, request, err)
	}
}

//Function for handling the calls to the default page or index
func defaultHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("Entering defaultHandler")

	//get URL.Path
	url := request.URL.Path[1:]
	//debug text
	utility.WriteInfo("URL = " + url)

	//check for time
	if url == "/" || url == "index/" || url == "index" || url == "" {

		name := auth.GetName(request)
		if name == "" {
			http.Redirect(response, request, "http://localhost:8080/login", http.StatusFound)
		} else {
			//debug text
			utility.WriteTrace("Displaying greeting")

			//create data structure
			data := view.TimeData{
				name, "", "",
			}

			//execute template
			err := view.ShowPage(response, "index", data)

			//process error
			if err != nil {
				errorHandler(response, request, err)
			}
		}

	} else {
		//send 404
		http.NotFound(response, request)
	}
}

//Function for handling calls to the login page
func loginHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("Entering loginHandler")

	//in case you went here directly without logging out
	name := auth.GetName(request)
	if name != "" {
		//redirect to home page
		http.Redirect(response, request, "http://localhost:8080/index", http.StatusFound)
	} else {

		//parse query from URL
		v, err := url.ParseQuery(request.URL.RawQuery)
		name := v.Get("name")

		//print any errors
		if err != nil {
			//send error to page
			errorHandler(response, request, err)
		} else if len(name) != 0 { //Check if there is text to parse

			//debug text
			utility.WriteInfo("len(name) != 0")

			//create and set authentication information
			err = auth.SetAuthInfo(response, name)

			//process error
			if err != nil {
				errorHandler(response, request, err)
			} else {
				//debug text
				utility.WriteInfo("Cookie set")

				//redirect to home page
				http.Redirect(response, request, "http://localhost:8080/index", http.StatusFound)
			}
		} else {
			//debug text
			utility.WriteTrace("displaying login form")

			//get message to display
			var message string
			if len(v) == 0 { //no query text
				message = "What is your name, Earthling?"
			} else {
				message = "C'mon, I need a name."
			}

			//create data structure
			data := view.TimeData{
				"", message, "",
			}

			//execute template
			err := view.ShowPage(response, "login", data)

			//process error
			if err != nil {
				errorHandler(response, request, err)
			}
		}
	}
}

//Function for handling logouts
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	//print debug
	utility.WriteTrace("Entering logoutHandler")

	//remove authentication information
	auth.DelAuthInfo(response)

	//create data structure
	data := view.TimeData{
		"", "", "",
	}

	//the message "Good-bye." is displayed for 10 seconds
	//execute template
	err := view.ShowPage(response, "logout", data)

	//process error
	if err != nil {
		errorHandler(response, request, err)
	}
}

//Handler for the web page.  One handler for all pages, URL.Path is used for sub pages.
func timeHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("Entering timeHandler")

	//get time and name
	time := time.Now().Format("3:04:05 PM")
	name := auth.GetName(request)

	//check if name was specified
	if name != "" {
		name = ", " + name
	}

	//create data structure
	data := view.TimeData{
		name, "", time,
	}

	//execute template
	err := view.ShowPage(response, "time", data)

	//process error
	if err != nil {
		errorHandler(response, request, err)
	}
}

//Function to handle errors
func errorHandler(response http.ResponseWriter, request *http.Request, err error) {
	//debug text
	view.ShowError(response, err)

	//debug text
	utility.WriteCritical(err.Error())
}
