package main

import (
	"log"
	"net/http"
	"time"
)

type XavierResponse struct {
	resp    *http.Response
	label   string
	latency time.Duration
}

type XavierConf struct {
	serviceList map[string]ServiceConf
	timeout     time.Duration
}

type ServiceConf struct {
    url string
}

func serviceMonitor(conf *XavierConf, responseStream chan<- XavierResponse) {
    Client := &http.Client{}

    for label, config := range conf.serviceList {
        url := config.url

        request, err := http.NewRequest("HEAD", url, nil)
        if err != nil {
            log.Fatalln(err)
        }
        request.Header.Set("User-Agent", "Xavier monitoring spider v0.1")

        startTime := time.Now()
        response, err := Client.Do(request)
        if err != nil {
            log.Fatalln(err)
        }
        endTime := time.Now()
        totalTime := endTime.Sub(startTime)

        monitorResponse := XavierResponse{
            response,
            label,
            totalTime,
        }
        responseStream <- monitorResponse

    }
    close(responseStream)
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

    go serviceMonitor(testConf, responseChan)

    testResponse := <-responseChan

    log.Println("Printing responses:\n", testResponse)

    for response := range responseChan {
        log.Println(response)
    }
}
