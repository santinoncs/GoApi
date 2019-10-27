/**

Create a server listening at port 8080, with the following API:

Response format: JSON like {"status":<status_code>,"body":"<response_body>"}
For example, successful response: {"status":"OK", "body":"Hello World"}
             error response: {"status":"Error", "body":"Method not supported"}

HTTP GET:
/api - returns the list of supported methods
/api/ping - returns the string "pong"
/api/set?key=<key>&value=<value> - sets a string value to a given string key
/api/get?key=<key> - returns the value saved at the key or an empty string if the value has not been set
/api/rpush?key=<key>&value=<value> - adds a value to the end of a list
/api/rpop?key=<key> - returns the value from the end of the list
/api/llen?key=<key> - returns the length of the list


any other url - returns the string "Method not supported"
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"db/db"
)

// Response : here you tell us what Response is
type Response struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

var di Db
var li ListItem

func main() {
	di = Db{}
	

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)
}

func response(status string, body string) {

}

func handler(w http.ResponseWriter, r *http.Request) {

	var value string
	var status string

	// call needed API method - first 4 methods

	if r.URL.Path == "/api" {
		value = "The methods are: /api /api/ping /api/set?key=<key>&value=<value> /api/get?key=<key>"
		status = "OK"
	}

	// call ping url

	if r.URL.Path == "/api/ping" {
		value = "Pong"
		status = "OK"
	}

	// any other url - returns the string "Method not supported"

	m1, _ := regexp.MatchString("api/(ping|set|get|rpush|rpop|llen).*", r.URL.Path)
	m2, _ := regexp.MatchString("/api$", r.URL.Path)

	if !m1 && !m2 {
		value = "Method not supported"
		status = "Error"
	}

	// implement URL parsing

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}

	values := make(map[string]string)

	for key := range params {
		values[key] = params.Get(key)
	}

	si := StringItem{}

	if r.URL.Path == "/api/get" {
		s, err := di.get(values["key"])
		if err == nil {
			stringItem := StringItem{}
			if s.dataType() == "string" {
				stringItem, _ = s.(StringItem)
			}
			value = stringItem.get()
			status = "OK"

		} else {
			value = "The key was not defined"
			status = "Error"
		}


	}

	if r.URL.Path == "/api/set" {
		si.set(values["value"])
		value = values["key"]
		di.set(value, si)
		status = "OK"

	}

	if r.URL.Path == "/api/rpush" {

		li,err := di.get(values["key"])
		if err == nil {
			fmt.Println("ya existe el array")
			listItem := ListItem{}
			// how do you know this type is listitem, here?
			if li.dataType() == "list" {
				listItem, _ = li.(ListItem)
			
				listItem.value = append(listItem.value, values["value"])

				di.set(values["key"], listItem)
				value = values["value"]
				status = "OK"
			}

		} else {
			fmt.Println("No existe el array")
			listItem := ListItem{}
			listItem, _ = li.(ListItem)
			listItem.value = append(listItem.value, values["value"])
			value = values["value"]
			di.set(values["key"], listItem)
			status = "OK"

		}
		
		

	}

	if r.URL.Path == "/api/rpop" {
		l, err := di.get(values["key"])
		if err == nil {
			listItem := ListItem{}
			if l.dataType() == "list" {
				listItem, _ = l.(ListItem)
			
				val := listItem.value
				fmt.Println(listItem.value)
				
				if len(listItem.value) > 0 {
					
					value = val[len(val)-1]
					listItem.value[len(val)-1] = ""
					listItem.value = listItem.value[:len(listItem.value)-1]
					fmt.Println(listItem.value)
					di.set(values["key"], listItem)
					status = "OK"

				} else {
					value = "The slice is empty"
					status = "OK"
				}

			}

		} else {
			value = "The key was not defined"
			status = "Error"

		}


	}

	if r.URL.Path == "/api/llen" {
		l, err := di.get(values["key"])

		if err == nil {
			listItem := ListItem{}
			if l.dataType() == "list" {
				listItem, _ = l.(ListItem)
			}
			val := listItem.value
			value = strconv.Itoa(len(val))
			status = "OK"

		} else {
			value = "there is no list"
			status = "Error"
		}


	}

	response := Response{Status: status, Body: value}
	responseJSON, _ := json.Marshal(response)
	fmt.Fprintf(w, "Response: %s\n", responseJSON)

}
