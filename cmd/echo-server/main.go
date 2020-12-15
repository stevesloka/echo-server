package main

import (
	"flag"
	"fmt"
	htmltemplate "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	echoText      string
	responseDelay time.Duration
	listenPort    int
	delay         <-chan time.Time
)

const (
	templatesBase = "templates"
)

// echoInfo is used to store dynamic properties on
// the echo template
type echoInfo struct {
	App     string
	Host    string
	Request string
	Headers http.Header

	BackgroundColor string
}

func getRequest(w http.ResponseWriter, r *http.Request) {

	// Add delay if enabled
	if responseDelay > 0 {
		<-delay
	}

	name, err := os.Hostname()
	if err != nil {
		fmt.Println("err: ", err)
	}

	outputText := fmt.Sprintf("ECHO Request Server: \n--------------------\n")
	outputText += fmt.Sprintf("App: \n    %s\n", echoText)
	outputText += fmt.Sprintf("Host: \n    %s\n", name)

	headers := r.Header
	outputText += fmt.Sprintf("Request: \n    http://%s%s\n", r.Host, r.RequestURI)
	outputText += fmt.Sprintf("Headers: \n    %s\n", headers)

	backgroundColor := "bg-light"

	if val := r.Header.Get("iscanary"); val == "true" {
		backgroundColor = "bg-primary"
	}

	data := &echoInfo{
		App:             echoText,
		Host:            name,
		Request:         fmt.Sprintf("http://%s%s\n", r.Host, r.RequestURI),
		Headers:         r.Header,
		BackgroundColor: backgroundColor,
	}

	if r.URL.Query().Get("format") == "text" {
		w.Write([]byte(outputText))
	} else {
		serveTemplate("echo.tmpl", data, w)
	}

	// Log to stdout
	fmt.Println(outputText)
}

func serveTemplate(tmplFile string, data interface{}, w http.ResponseWriter) {
	var (
		templatePath string
		templateData []byte
		err          error
	)

	templatePath = filepath.Join(templatesBase, tmplFile)
	templateData, err = ioutil.ReadFile(templatePath)

	if err != nil {
		log.Errorf("Failed to find template asset: %s at path: %s", tmplFile, templatePath)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := htmltemplate.New(tmplFile)
	tmpl, err = tmpl.Parse(string(templateData))
	if err != nil {
		log.Errorf("Failed to parse template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, tmplFile, data)
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
