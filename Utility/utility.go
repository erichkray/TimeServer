// Author: Erich Ray
// Title: Utility
// Version 3
// Created on 2/2/2015
// Utility methods for TimeServer

package utility

//import packages
import (
	"flag"
	"fmt"
	log "github.com/seelog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}
//global variable
var (
	verbose *bool
	port    *int
	Port = func() int {
		return *port
	}
	authPort *int
	AuthPort = func() int{
		return *authPort
	}
	debug   *bool
	logstr     *string
	logger log.LoggerInterface
)

func ParseFlags(version string){
	//parse flag
	verbose = flag.Bool("v", false, "a bool")
	port = flag.Int("port", 8080, "port for webserver")
	logstr = flag.String("log", "", "load location for seelog")	
	debug = flag.Bool("debug", false, "turn optional for debug spew")
	flag.Parse()
	
		//check if log was specified
	if *logstr == "" {
		WriteInfo("No log config xml loaded, outputing to console only.")
	} else {

		//load config file and check for errors
		err := SetLogConfig(*logstr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	
	if *verbose == true {
	fmt.Println(version)
		os.Exit(0)
	}
}

func SetLogConfig(logConfigPath string) error {
	logger, err := log.LoggerFromConfigAsFile(logConfigPath)

	if err != nil {
		return err
	}

	log.ReplaceLogger(logger)
	return nil
}

func WriteTrace(msg string) {
	printDebug(msg)
	log.Trace(msg)
}

func WriteInfo(msg string) {
	printDebug(msg)
	log.Info(msg)
}

func WriteCritical(msg string) {
	printDebug(msg)
	log.Critical(msg)
}

func FlushLog() {
	logger.Flush()
}

//Function to create UUID
func UUIDCreator() (string, error) {
	//debug text
	printDebug("uuidCreator")
	if runtime.GOOS == "windows" {
		out, err := exec.Command("C:\\Program Files (x86)\\microsoft sdks\\Windows\\v7.1A\\Bin\\x64\\Uuidgen.Exe").Output()
		return strings.TrimRight(string(out), "\r\n"), err
	} else {
		out, err := exec.Command("uuidgen").Output()
		return string(out), err
	}
}

//Function to print debug text to Stderr
func printDebug(msg string) {
	if *debug {
		fmt.Fprintf(os.Stderr, "Debug: %s\n", msg)
	}
}

//print usage for TimeServer
func usage() {
	fmt.Println("TimeServer Usage:")
	fmt.Println("Creates a web server for the url http://localhost:8080/\n")
	fmt.Println("\\TimeServer.go [single argument]\n")
	fmt.Println("Arguments:")
	fmt.Println("Port Number: Can be any integer specifying a port for the server.  Default is 8080.")
	fmt.Println("Verbose: Pass -v to display version number.")
}
