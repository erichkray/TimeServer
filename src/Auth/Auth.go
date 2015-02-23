// Author: Erich Ray
// Title: CookieHandler
// Version 1
// Created on 1/29/2015
// Handles Cookies

package auth

//import packages
import (
	"errors"
	"io/ioutil"
	"Utility"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

//global variable
var (
	session_uuid string
)

var jar = struct {
	sync.RWMutex
	cookies map[string]string
}{cookies: make(map[string]string)}

//Get the cookie from the request
func findCookie(request *http.Request, session_uuid string) (*http.Cookie, error) {

	//if there is no cookie or there a cookie but no corresponding name,
	//display the login form.
	if jar.cookies == nil {

		return nil, nil
	}

	cookie, err := request.Cookie("user")
	if err != nil {
		return nil, err
	} else if cookie.Value != session_uuid {
		err := errors.New("UUID did not match for this session.")
		return nil, err
	}

	return cookie, nil
}

//remove the cookie from the response
func removeCookie(response http.ResponseWriter, session_uuid string) {

	//RPC for name to clear map
	strStub := "/set?cookie=" + session_uuid + "&name="
	_, err := http.Get("http://localhost:" + strconv.Itoa(utility.AuthPort()) + strStub)

	if err != nil {
		utility.WriteCritical(err.Error())
	}

	//expire cookie
	exp_cookie := http.Cookie{Name: "user", Value: session_uuid, MaxAge: -1}
	http.SetCookie(response, &exp_cookie)
}

func createCookie(response http.ResponseWriter, session_uuid string, name string) *http.Cookie {
	//make cookie
	cookie := http.Cookie{Name: "user", Value: session_uuid, Path: "/"}

	//set cookie
	http.SetCookie(response, &cookie)
	addNameToServer(name, session_uuid)

	return &cookie
}

//add the name to the session map
func addNameToServer(name string, session_uuid string) {

	//debug text
	utility.WriteInfo("AddName(" + name + ", " + session_uuid + ")")

	//RPC for name to set map
	strStub := "/set?cookie=" + session_uuid + "&name=" + name
	utility.WriteInfo("Sending: http://localhost:" + strconv.Itoa(utility.AuthPort()) + strStub)
	resp, err := http.Get("http://localhost:" + strconv.Itoa(utility.AuthPort()) + strStub)

	if err != nil {
		utility.WriteCritical(err.Error())
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utility.WriteCritical(err.Error())
		} else if len(contents) > 0 {
			outStr := strings.TrimRight(string(contents), "\r\n")
			outStr = strings.TrimRight(string(outStr), "\n")
			utility.WriteInfo("Server Response: " + string(outStr))
		}
	}
}

//get the name from the session map
func getNameFromServer(session_uuid string) string {
	//debug text
	utility.WriteInfo("GetName(" + session_uuid + ")")

	//RPC for name to get from map
	strStub := "/get?cookie=" + session_uuid
	utility.WriteInfo("Sending: http://localhost:" + strconv.Itoa(utility.AuthPort()) + strStub)
	resp, err := http.Get("http://localhost:" + strconv.Itoa(utility.AuthPort()) + strStub)

	var outStr string
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utility.WriteCritical(err.Error())
	} else if len(contents) > 0 {
			outStr = strings.TrimRight(string(contents), "\r\n")
			outStr = strings.TrimRight(string(outStr), "\n")
			utility.WriteInfo("Server Response: " + string(outStr))
		}
	return string(contents)
}

func GetName(request *http.Request) string {
	utility.WriteInfo("Auth.GetName()")
	utility.WriteInfo("session_uuid: " + session_uuid)

	//check session_uuid
	if session_uuid == "" {
		//return empty string, will trigger auth setting
		utility.WriteInfo("session_uuid empty, returning")
		return ""
	} else {

		//try get cookie
		cookie, _ := findCookie(request, session_uuid)

		//validate cookie
		if cookie == nil {
			utility.WriteInfo("Cookie = nil")
			return ""
		} else {

			//get name from cookie
			//cookiehandler.GetName returns empty string on error
			name := getNameFromServer(session_uuid)
			return name
		}
	}
}

func SetAuthInfo(response http.ResponseWriter, name string) error {
	utility.WriteInfo("Auth.SetAuthInfo()")

	//set UUID
	var err error
	session_uuid, err = utility.UUIDCreator()

	//check for errors
	if err != nil {
		return err
	}

	//create and set cookie
	createCookie(response, session_uuid, name)

	//placeholder for more error handling
	return nil
}

func DelAuthInfo(response http.ResponseWriter) error {
	//remove cookie
	removeCookie(response, session_uuid)

	//reset session_uuid
	session_uuid = ""

	//placeholder for more error handling
	return nil
}
