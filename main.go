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
	serviceList *map[string]string
	timeout     int8
}

func serviceMonitor(label string) XavierResponse {
	client := &http.Client{}
	request, err := http.NewRequest("HEAD", "http://www.google.com", nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("User-Agent", "Xavier monitoring spider v0.1")

	startTime := time.Now()
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	return XavierResponse{
		response,
		label,
		totalTime,
	}
}

func main() {
	response := serviceMonitor("Google")

	log.Println(response.resp.StatusCode, ":", response.label, ":", response.latency, "ms")
}
