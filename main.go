package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var originalHash uint32
var lastResponseBody string
var url string
var interval int
var config Config

func init() {
	flag.StringVar(&url, "u", "", "URL to target.")
	flag.IntVar(&interval, "n", 5, "Time interval (seconds) between requests.")
	flag.Parse()

	file, err := os.Open("config.json")
	if err != nil {
		die("Failed to open 'config.json' - %s", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			die("Failed to close 'config.json' file - %s", err)
		}
	}()

	decode := json.NewDecoder(file)
	err = decode.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	if url == "" {
		die("You must provide a URL. Try 'snoop --help'")
	}

	check(url)

	fmt.Printf("Monitoring '%v' every %v seconds.\n", url, interval)
	for {
		calculatedHash := hashURL(url)

		if calculatedHash != originalHash {
			notify(url)
			os.Exit(0)
		}

		time.Sleep(time.Duration(interval) * time.Second)

	}

}
