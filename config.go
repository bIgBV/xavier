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
}

// Service type stores configuration regarding indidual services to be
// monitored. This is data such as the URL, auth, timeout, etc.
type Service struct {
	URL     string
	Timeout time.Duration
}

func parseConfig(fname string) Conf {
	file, err := ioutil.ReadFile(fname)

	if err != nil {
		log.Fatal("Unable to read file: ", fname)
		panic(err)
	}

	var conf Conf

	if err := toml.Unmarshal(file, &conf); err != nil {
		log.Fatal("Configuration parse error: \n", err)
		panic(err)
	}

	log.Println("Configuration \n", conf)
	return conf
}
