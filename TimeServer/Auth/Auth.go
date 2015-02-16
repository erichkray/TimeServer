// Author: Erich Ray
// Title: Auth
// Version 1
// Created on 2/15/2015
// Authority package for the TimeServer project

package auth

//import packages
import (
	"local/TimeServer/CookieHandler"
	"local/TimeServer/Utility"
	"net/http"
)

//global variable
var (
	session_uuid string
	debug        *bool
)

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
		cookie, _ := cookiehandler.FindCookie(request, session_uuid)

		//validate cookie
		if cookie == nil {
			return ""
		} else {

			//get name from cookie
			//cookiehandler.GetName returns empty string on error
			name := cookiehandler.GetName(session_uuid)
			utility.WriteInfo("name: " + name)
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
	cookiehandler.CreateCookie(response, session_uuid)
	cookiehandler.AddName(name, session_uuid)

	//placeholder for more error handling
	return nil
}

func DelAuthInfo(response http.ResponseWriter) error {
	//remove cookie
	cookiehandler.RemoveCookie(response, session_uuid)

	//reset session_uuid
	session_uuid = ""

	//placeholder for more error handling
	return nil
}
