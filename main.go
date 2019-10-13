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
	li = ListItem{}

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)
}

func handler(w http.ResponseWriter, r *http.Request) {

	// call needed API method - first 4 methods

	if r.URL.Path == "/api" {
		fmt.Fprintf(w, "The methods are: \n/api\n/api/ping\n/api/set?key=<key>&value=<value>\n/api/get?key=<key>\n")
	}

	// call ping url

	if r.URL.Path == "/api/ping" {
		response := Response{Status: "OK", Body: "Pong"}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)
	}

	// any other url - returns the string "Method not supported"

	m1, _ := regexp.MatchString("api/(ping|set|get|rpush|rpop|llen).*", r.URL.Path)
	m2, _ := regexp.MatchString("/api$", r.URL.Path)

	if !m1 && !m2 {
		response := Response{Status: "Error", Body: "Method not supported"}
		responseJSON, _ := json.Marshal(response)

		fmt.Fprintf(w, "Response: %s\n", responseJSON)
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
		var value string
		if err == nil {
			stringItem := StringItem{}
			if s.dataType() == "string" {
				stringItem, _ = s.(StringItem)
			}
			value = stringItem.get()

		} else {
			value = "The key was not defined"
		}

		response := Response{Status: "OK", Body: value}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/api/set" {
		si.set(values["value"])
		st := values["key"]
		di.set(st, si)
		response := Response{Status: "OK written", Body: st}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)
	}

	if r.URL.Path == "/api/rpush" {
		li.rpush(values["value"])
		di.set(values["key"], li)
		response := Response{Status: "OK written", Body: values["value"]}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)
	}

	if r.URL.Path == "/api/rpop" {
		s, err := di.get(values["key"])

		var output string
		if err == nil {
			listItem := ListItem{}
			if s.dataType() == "list" {
				listItem, _ = s.(ListItem)
			}
			val := listItem.get()
			output = val[len(val)-1]

		} else {
			output = "The key was not defined"

		}

		response := Response{Status: "OK written", Body: output}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)
	}

	if r.URL.Path == "/api/llen" {
		s, err := di.get(values["key"])

		var output string
		if err == nil {
			listItem := ListItem{}
			if s.dataType() == "list" {
				listItem, _ = s.(ListItem)
			}
			val := listItem.get()
			output = strconv.Itoa(len(val))

		} else {
			output = "there is no list"
		}

		response := Response{Status: "OK", Body: output}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)
	}

}
