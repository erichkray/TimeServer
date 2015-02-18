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
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//global variables
var (

	//private vars
	verbose *bool
	port    *int
	authPort *int
	authTimeout *int		
	avgResponse *int	
	deviation *int
	debug  *bool
	logstr *string
	logger log.LoggerInterface
	maxInFlight *int
	
	//public vars (Accessors)
	Port = func() int {
		return *port
	}	
	AuthPort = func() int {
		return *authPort
	}
	AuthTimeout = func() int {
		return *authTimeout
	}
	AvgResponse = func() int {
		return *avgResponse
	}
	Deviation = func() int {
		return *deviation
	}
	Delay = func() time.Duration {
		return time.Duration(rand.Int() + *deviation)
	}
	MaxInFlight = func() int {
		return *maxInFlight
	}
)
 
//PUBLIC METHODS

//initialize the utility settings
func Init(version string) {
	//parse flag
	verbose = flag.Bool("v", false, "Enable/Disable verbose messaging (default: false)")
	port = flag.Int("port", 8080, "Port for webserver (default: 8080)")
	logstr = flag.String("log", "", "Load location for seelog (default: NA, print to screen only)")
	debug = flag.Bool("debug", false, "Turn optional debug spew (default: false)")
	authPort = flag.Int("authport", 9090, "Port for authserver (default: 9090)")
	authTimeout = flag.Int("authtimeout-ms", 500, "Timeout to terminate auth request (default: 500)")
	avgResponse = flag.Int("avg-response-ms", 100, "Seed for random number generator (default: 100)")
	deviation = flag.Int("deviation-ms", 100, "Deviation value for random number generator (default: 100)")
	maxInFlight = flag.Int("max-inflight", 0, "Maximum number of in-flight time requests the server can handle")
	flag.Parse()

	//create random
	rand.Seed(int64(*avgResponse))
	
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

//set the log config file
func SetLogConfig(logConfigPath string) error {
	logger, err := log.LoggerFromConfigAsFile(logConfigPath)

	if err != nil {
		return err
	}

	log.ReplaceLogger(logger)
	return nil
}

//write trace helper
func WriteTrace(msg string) {
	printDebug(msg)
	log.Trace(msg)
}

//write info helper
func WriteInfo(msg string) {
	printDebug(msg)
	log.Info(msg)
}

//write critical helper
func WriteCritical(msg string) {
	printDebug(msg)
	log.Critical(msg)
}

//flush the log
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

//PRIVATE METHODS

//Function to print debug text to Stderr
func printDebug(msg string) {
	if *debug {
		fmt.Fprintf(os.Stderr, "Debug: %s\n", msg)
	}
}
