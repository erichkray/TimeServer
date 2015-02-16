// Author: Erich Ray
// Title: AuthServer
// Version 1
// Created on 1/15/2015
// Creates an auth server for the timeserver project
package main

//import packages
import (
	"flag"
	"fmt"
	"local/TimeServer/Utility"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const version = "Version: 1"

//global variable
var (
	verbose *bool
	port    *int
	debug   *bool
	log     *string
)

//main function for Time Server
func main() {
	//parse command flags
	utility.ParseFlags(version)

	//defer a flush for on quit
	defer utility.FlushLog()

	//print debug
	utility.WriteTrace("Starting Main")
	
	//write version number to log
	utility.WriteInfo(version)

	//setup handler for web page requests
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/set/", setHandler)

	//set listen and serve for port, checking for error
	err := http.ListenAndServe(":"+strconv.Itoa(utility.Port()), nil)

	//if error was returned, send error to standard output
	if err != nil {
		utility.WriteCritical(err.Error())
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func getHandler(){

}

func setHandler(){

}
