// Author: Erich Ray
// Title: CookieHandler
// Version 1
// Created on 1/29/2015
// Handles Cookies

package cookiehandler

//import packages
import (
	"errors"
	"net/http"
	"sync"
)

var jar = struct {
	sync.RWMutex
	cookies map[string]string
}{cookies: make(map[string]string)}

//Abstraction for finding a cookie
func FindCookie(request *http.Request, session_uuid string) (*http.Cookie, error) {

	//if there is no cookie or there a cookie but no corresponding name,
	//display the login form.
	if jar.cookies == nil {

		return nil, nil
	}

	/*cookieStr := "Cookies: "
	cookie_req := request.Cookies()
	for _, c := range cookie_req {
		cookieStr += c.Name + " "
	}*/

	cookie, err := request.Cookie("user")
	if err != nil {
		return nil, err
	} else if cookie.Value != session_uuid {
		err := errors.New("UUID did not match for this session.")
		return nil, err
	}

	return cookie, nil
}

//Abstraction for removing a cookie
func RemoveCookie(response http.ResponseWriter, session_uuid string) {

	//delete name from map
	jar.Lock()
	jar.cookies[session_uuid] = ""
	jar.Unlock()

	//expire cookie
	exp_cookie := http.Cookie{Name: "user", Value: session_uuid, MaxAge: -1}
	http.SetCookie(response, &exp_cookie)
}

func AddName(name string, session_uuid string) {
	//add name to the cookie map
	jar.Lock()
	jar.cookies[session_uuid] = name
	jar.Unlock()
}

func GetName(session_uuid string) string {
	//get name from the cookie map
	jar.RLock()
	name := jar.cookies[session_uuid]
	jar.RUnlock()

	return name
}

//Abstraction for making Jar
func MakeJar() {
	//make cookie map
	jar.Lock()
	jar.cookies = make(map[string]string)
	jar.Unlock()
}

func CreateCookie(response http.ResponseWriter, session_uuid string) *http.Cookie {
	//make cookie
	cookie := http.Cookie{Name: "user", Value: session_uuid, Path: "/"}

	//set cookie
	http.SetCookie(response, &cookie)
	return &cookie
}
