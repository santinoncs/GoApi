# GoApi




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
