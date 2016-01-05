// Package Xavier provides a simple tool to monitor various services.
package main

import (
	"log"
	"net/http"
	"time"
)

// XavierResponse type is what is sent to persistance adapters.
// Since the adapters can be configured to persist whatever data from
// the response they choose to, it contains a reference to the response
// for a particular request.
type XavierResponse struct {
	resp    *http.Response
	label   string
	latency time.Duration
}

// XavierConf type is the single source of truth for configuration of
// the system.
type XavierConf struct {
	serviceList map[string]ServiceConf
	timeout     time.Duration
}

// ServiceConf type stores configuration regarding indidual services to be
// monitored. This is data such as the URL, auth, timeout, etc.
type ServiceConf struct {
    url string
}


func main() {

    var serviceListConf = map[string]ServiceConf{
        "Github": ServiceConf{url: "http://github.com"},
        "Rbox": ServiceConf{url: "http://www.recruiterbox.com"},
        "Google": ServiceConf{url: "http://www.google.com"},
        "Reddit": ServiceConf{url: "http://www.reddit.com"},
    }

    testConf := &XavierConf {
        serviceList: serviceListConf,
        timeout: time.Second * 10,
    }

    var responseChan = make(chan XavierResponse)

    go MonitorService(testConf, responseChan)

    testResponse := <-responseChan

    log.Println("Printing responses:\n", testResponse)

    for response := range responseChan {
        log.Println(response)
    }
}
