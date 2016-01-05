package main

import (
    "time"
    "net/http"
    "log"
)

func MonitorService(conf *XavierConf, responseStream chan<- XavierResponse) {
    Client := &http.Client{}
    for {
        for label, config := range conf.serviceList {
            go func(label string, config ServiceConf) {
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
            }(label, config)

        }
        time.Sleep(4)
    }
    close(responseStream)
}
