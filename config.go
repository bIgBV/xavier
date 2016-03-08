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
	serviceList map[string]Service `yaml:"services,flow"`
	timeout     time.Duration      `yaml:"timeout"`
}

// Service type stores configuration regarding indidual services to be
// monitored. This is data such as the URL, auth, timeout, etc.
type Service struct {
	url string `yaml:",flow"`
}

func parseConfig(fname string) (conf *Conf) {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Unable to read file: ", fname)
		panic(err)
	}

	err = yaml.Unmarshal(file, &conf)

	return conf
}
