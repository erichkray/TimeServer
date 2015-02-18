// Author: Erich Ray
// Title: AuthServer
// Version 1
// Created on 1/15/2015
// Creates an auth server for the timeserver project
package main

//import packages
import (
	"fmt"
	"local/TimeServer/Utility"
	"net/http"
	//"net/url"
	"os"
	"strconv"
	"sync"
)

const version = "Version: 1"

//global variable
var (
	verbose *bool
	port    *int
	debug   *bool
	log     *string
	jar     cookieJar
)

type cookieJar struct {
	sync.RWMutex
	cookies map[string]string
}

//main function for Time Server
func main() {
	//parse command flags
	utility.Init(version)

	//defer a flush for on quit
	defer utility.FlushLog()

	//print debug
	utility.WriteTrace("Starting Main")

	//write version number to log
	utility.WriteInfo(version)

	//create cookieJar
	jar.cookies = make(map[string]string)

	//setup handler for web page requests
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/set/", setHandler)

	//set listen and serve for port, checking for error
	err := http.ListenAndServe(":"+strconv.Itoa(utility.AuthPort()), nil)

	//if error was returned, send error to standard output
	if err != nil {
		utility.WriteCritical(err.Error())
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

//Function for handling the calls to the default page or index
func defaultHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("defaultHandler()")

	//debug text
	utility.WriteInfo("returning 404")
	
	//return 404
	http.NotFound(response, request)
}

func getHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("getHandler()")
	utility.WriteInfo(request.URL.RawQuery)

	//get uuid
	err := request.ParseForm()

	//vars, err := url.ParseQuery(request.URL.RawQuery)
	//session_uuid := vars.Get("cookie")

	//utility.WriteInfo("session_uuid: " + session_uuid)

	if err != nil {
		fmt.Fprintln(response, err.Error())
		utility.WriteCritical(err.Error())
	} else {
		//get name from the cookie map
		//vars := request.Form
		session_uuid := request.PostFormValue("cookie")
		utility.WriteInfo("session_uuid: " + session_uuid)

		name := jar.Get(session_uuid)

		if name == "" {
			//print name to response
			fmt.Fprintln(response, "No Name!")
		} else {
			//print name to response
			fmt.Fprintln(response, name)
		}
	}
}

func setHandler(response http.ResponseWriter, request *http.Request) {
	//debug text
	utility.WriteTrace("setHandler()")
	utility.WriteInfo(request.URL.RawQuery)

	//get uuid
	err := request.ParseForm()

	//vars, err := url.ParseQuery(request.URL.RawQuery)
	//session_uuid := vars.Get("cookie")
	//name := vars.Get("name")

	//utility.WriteInfo("session_uuid: " + session_uuid + " - name: " + name)

	if err != nil {
		fmt.Fprintln(response, err.Error())
		utility.WriteCritical(err.Error())
	} else {
		vars := request.Form
		session_uuid := vars.Get("cookie")
		name := vars.Get("name")

		utility.WriteInfo("session_uuid: " + session_uuid + " - name: " + name)

		jar.Add(session_uuid, name)

		//print name to response
		fmt.Fprintln(response, "cookie set")
	}
}

//add name to the cookie map
func (cj *cookieJar) Add(session_uuid string, name string) {
	cj.Lock()
	cj.cookies[session_uuid] = name
	cj.Unlock()
}

//get the name from the cookie map
func (cj *cookieJar) Get(session_uuid string) string {
	cj.RLock()
	name := cj.cookies[session_uuid]
	cj.RUnlock()
	return name
}
