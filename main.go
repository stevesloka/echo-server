package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

var (
	echoText      string
	responseDelay time.Duration
	listenPort    int
	delay         <-chan time.Time
)

func getRequest(w http.ResponseWriter, r *http.Request) {

	// Add delay if enabled
	if responseDelay > 0 {
		<-delay
	}

	w.Write([]byte(fmt.Sprintf("ECHO Request Server: \n--------------------\n")))
	w.Write([]byte(fmt.Sprintf("App: \n    %s\n", echoText)))

	headers := r.Header
	w.Write([]byte(fmt.Sprintf("Request: \n    http://%s%s\n", r.Host, r.RequestURI)))
	w.Write([]byte(fmt.Sprintf("Headers: \n    %s\n", headers)))
}

func init() {
	flag.StringVar(&echoText, "echotext", "", "enter text to echo back to the user")
	flag.DurationVar(&responseDelay, "response-delay", 0, "")
	flag.IntVar(&listenPort, "listen-port", 8080, "The port used to listen on. Defaults to 8080")
}

func main() {
	flag.Parse()
	delay = time.Tick(responseDelay)
	http.HandleFunc("/", getRequest)

	fmt.Printf("Server started! Listening on port %q", fmt.Sprintf(":%d", listenPort))
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
