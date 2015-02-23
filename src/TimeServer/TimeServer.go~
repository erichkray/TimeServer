// Author: Erich Ray
// Title: TimeServer
// Version 4
// Created on 1/15/2015
// Creates a web server for the url http://localhost:<port>/
package main

//import packages
import (
	//"flag"
	"fmt"
	"Auth"
	"View"
	"Utility"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

const version = "Version: 4"

var (
	timeRequests TimeRequestCounter
)

type TimeRequestCounter struct{
	sync.Mutex
	requests int
}

//main function for Time Server
func main() {

	//parse command flags
	utility.Init(version)

	//build template tree
	view.CreateSite()

	//defer a flush for on quit
	defer utility.FlushLog()

	//print debug
	utility.WriteTrace("Starting Main")

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
	err := http.ListenAndServe(":"+strconv.Itoa(utility.Port()), nil)

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
		"", "", "", strconv.Itoa(utility.Port()),
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
			http.Redirect(response, request, "http://localhost:"+strconv.Itoa(utility.Port())+"/login", http.StatusFound)
		} else {
			//debug text
			utility.WriteTrace("Displaying greeting")

			//create data structure
			data := view.TimeData{
				name, "", "", strconv.Itoa(utility.Port()),
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
	utility.WriteTrace("Checking if user went directly to login")
	name := auth.GetName(request)
	if name != "" {
		//redirect to home page
		http.Redirect(response, request, "http://localhost:"+strconv.Itoa(utility.Port())+"/index", http.StatusFound)
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
				//redirect to home page
				http.Redirect(response, request, "http://localhost:"+strconv.Itoa(utility.Port())+"/index", http.StatusFound)
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
				"", message, "", strconv.Itoa(utility.Port()),
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
		"", "", "", strconv.Itoa(utility.Port()),
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
	utility.WriteTrace("Incrementing timeRequests")
	
	//increment counter
	timeRequests.increment()
	
	//check if max number of concurrency has occurred
	if utility.MaxInFlight() != 0 && timeRequests.requests > utility.MaxInFlight() {
		http.Error(response, "Internal Server Error", 500)
	} else {
	
		//get time and name
		currentTime := time.Now().Format("3:04:05 PM")
		name := auth.GetName(request)

		//check if name was specified
		if name != "" {
			name = ", " + name
		}

		//create data structure
		data := view.TimeData{
			name, "", currentTime, strconv.Itoa(utility.Port()),
		}
		
		time.Sleep(utility.Delay() * time.Millisecond)

		//execute template
		err := view.ShowPage(response, "time", data)

		//process error
		if err != nil {
			errorHandler(response, request, err)
		}
	}
	
	//debug text
	utility.WriteTrace("Decrementing timeRequests")
	//decrement counter
	timeRequests.decrement()
}

//Function to handle errors
func errorHandler(response http.ResponseWriter, request *http.Request, err error) {
	//debug text
	view.ShowError(response, err)

	//debug text
	utility.WriteCritical(err.Error())
}

//increment counter
func (counter *TimeRequestCounter) increment() {
	counter.Lock()
	counter.requests++
	counter.Unlock()
}

//decrement counter
func (counter *TimeRequestCounter) decrement() {
	counter.Lock()
	counter.requests--
	counter.Unlock()
}
