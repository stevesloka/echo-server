package main

import (
	"flag"
	"fmt"
	"net/http"
)

	func getRequest(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("ECHO Request Server: \n--------------------\n")))
		w.Write([]byte(fmt.Sprintf("App: \n    %s\n", echoText)))

		headers := r.Header
		w.Write([]byte(fmt.Sprintf("Request: \n    http://%s%s\n", r.Host, r.RequestURI)))
		w.Write([]byte(fmt.Sprintf("Headers: \n    %s\n", headers)))
}

var echoText string
func init() {
	flag.StringVar(&echoText, "echotext", "", "enter text to echo back to the user")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", getRequest)
	http.ListenAndServe(":8080", nil)
}
