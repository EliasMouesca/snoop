package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/smtp"
	"os"
)

func hashURL(url string) (hash uint32) {
	resp, err := http.Get(url)
	if err != nil {
		die("Unable to request '%v' - %v", url, err.Error())
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			die("Unable to close response body - %v", err.Error())
		}
	}(resp.Body)

	h := fnv.New32a()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		die("There was a problem reading the response body", err.Error())
	}

	lastResponseBody = string(body)

	_, err = h.Write(body)
	if err != nil {
		die("There was a problem calculating the hash of '%v' - ", body, err.Error())
	}

	return h.Sum32()

}

func check(url string) {
	const NumberOfChecks = 3
	hashes := make([]uint32, NumberOfChecks)
	hashes[0] = hashURL(url)
	for i := 1; i < NumberOfChecks; i++ {
		hashes[i] = hashURL(url)
		if hashes[i] != hashes[i-1] {
			die("Received different responses for the first 3 requests (page may not be observable)")
		}
	}

	originalHash = hashes[0]

}

func notify(url string) {
	// Log it
	fmt.Println("The page changed!")
	fmt.Println("Received:")
	fmt.Println(lastResponseBody)

	// Send email
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	msg := []byte("Subject: Change on: '" + url + "'\r\n" +
		"\r\n" +
		"Detected a change on the page '" + url + "'\r\n")
	err := smtp.SendMail(config.Host+":"+config.Port, auth, config.From, []string{config.To}, msg)
	if err != nil {
		die("Could not send email - %v", err.Error())
	}
}

func die(format string, a ...any) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		print("Failed to print error!")
	}
	os.Exit(1)
}
