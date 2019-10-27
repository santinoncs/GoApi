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
	db "github.com/santinoncs/GoApi/db"
)

// Response : here you tell us what Response is
type Response struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

var di db.Db
var li db.ListItem

func main() {
	di = db.Db{}
	

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

	si := db.StringItem{}

	if r.URL.Path == "/api/get" {
		s, err := di.Get(values["key"])
		if err == nil {
			stringItem := db.StringItem{}
			if s.DataType() == "string" {
				stringItem, _ = s.(db.StringItem)
			}
			value = stringItem.Get()
			status = "OK"

		} else {
			value = "The key was not defined"
			status = "Error"
		}


	}

	if r.URL.Path == "/api/set" {
		si.Set(values["value"])
		value = values["key"]
		di.Set(value, si)
		status = "OK"

	}

	if r.URL.Path == "/api/rpush" {

		li,err := di.Get(values["key"])
		if err == nil {
			fmt.Println("ya existe el array")
			listItem := db.ListItem{}
			// how do you know this type is listitem, here?
			if li.DataType() == "list" {
				listItem, _ = li.(db.ListItem)
			
				listItem.Value = append(listItem.Value, values["value"])

				di.Set(values["key"], listItem)
				value = values["value"]
				status = "OK"
			}

		} else {
			fmt.Println("No existe el array")
			listItem := db.ListItem{}
			listItem, _ = li.(db.ListItem)
			listItem.Value = append(listItem.Value, values["value"])
			value = values["value"]
			di.Set(values["key"], listItem)
			status = "OK"

		}
		
		

	}

	if r.URL.Path == "/api/rpop" {
		l, err := di.Get(values["key"])
		if err == nil {
			listItem := db.ListItem{}
			if l.DataType() == "list" {
				listItem, _ = l.(db.ListItem)
			
				val := listItem.Value
				fmt.Println(listItem.Value)
				
				if len(listItem.Value) > 0 {
					
					value = val[len(val)-1]
					listItem.Value[len(val)-1] = ""
					listItem.Value = listItem.Value[:len(listItem.Value)-1]
					fmt.Println(listItem.Value)
					di.Set(values["key"], listItem)
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
		l, err := di.Get(values["key"])

		if err == nil {
			listItem := db.ListItem{}
			if l.DataType() == "list" {
				listItem, _ = l.(db.ListItem)
			}
			val := listItem.Value
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
