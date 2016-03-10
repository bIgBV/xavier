package main

import (
	"github.com/naoina/toml"
	"io/ioutil"
	"log"
	"time"
)

// Conf type is the single source of truth for configuration of
// the system.
type Conf struct {
	Services map[string]Service
	Timeout  time.Duration
}

// Service type stores configuration regarding indidual services to be
// monitored. This is data such as the URL, auth, timeout, etc.
type Service struct {
	URL string
}

func parseConfig(fname string) (conf *Conf) {
	file, err := ioutil.ReadFile(fname)

	if err != nil {
		log.Fatal("Unable to read file: ", fname)
		panic(err)
	}

	if err := toml.Unmarshal(file, &conf); err != nil {
		log.Fatal("Configuration parse error: \n", err)
		panic(err)
	}

	log.Println("Config %v\n", conf)
	return conf
}
