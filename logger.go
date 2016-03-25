// Package main provides helper functions for logging data
package main

import (
	"io"
	"log"
	"os"
)

func createLogFile(fname string) (f *os.File) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error opening log file")
	}
	return f
}

func LogData(data Persistance) {
	log.SetOutput(io.MultiWriter(os.Stderr, createLogFile("xavier.log")))

}
