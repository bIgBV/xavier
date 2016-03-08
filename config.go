package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

// Conf type is the single source of truth for configuration of
// the system.
type Conf struct {
	Services   []Service       `yaml:"services,flow"`
	Timeout    time.Duration   `yaml:"timeout"`
}

// Service type stores configuration regarding indidual services to be
// monitored. This is data such as the URL, auth, timeout, etc.
type Service struct {
    Label string   `yaml:"label"` 
	URL string     `yaml:"URL"`
}

func parseConfig(fname string) (conf *Conf) {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Unable to read file: ", fname)
		panic(err)
	}

	err = yaml.Unmarshal(file, &conf)
    if err != nil {
        log.Fatal("Config parse error")
        panic(err)
    }
    
    log.Println("Config %v\n", conf)
	return conf
}
