package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)


//========================================Main Function Starts Here========================================//
func main() {

	var (
		requestURL string
		password   string
		parsedUrl  *url.URL
		err        error
	)

	//flag is a package that allows us to add command line flags to applications
	flag.StringVar(&requestURL, "url", "", "target URL")
	flag.StringVar(&password, "password", "", "password needed to access url")

	flag.Parse()

	if parsedUrl, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Validation error: URL is invalid: %s\n Usage: ./http-get -h\n ", err)
		flag.Usage() // This method allows us to print the help statement on errors
		os.Exit(1)
	}

	client := http.Client{}

	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: /http-get <url>\n")
		os.Exit(1)
	}

	res, err := makeRequest(client, parsedUrl.String())
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			log.Fatalf("Error: %s (HTTPCode: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
		}
		log.Fatalf("Error: %s\n", err)
	}

	if res == nil {
		log.Fatalln("No Response")
	}

	fmt.Printf("Response: %s", res.GetResponse())
}

//============================================Application Functions===========================================//

func makeRequest(client http.Client, requestURL string) (Response, error) {

	response, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("http Get error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid output (HTTP code %d): %s", response.StatusCode, body)
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      "No valid JSON returned",
		}
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page Unmarshal error: %s", err),
		}
	}
	switch page.Name {
	case "words":
		var words Words

		json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words Unmarshal error: %s", err),
			}
		}

		return words, nil

	case "occurrence":
		var occurence Occurrence

		json.Unmarshal(body, &occurence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurences Unmarshal error: %s", err),
			}
		}

		// Checking if the value exists
		if val, ok := occurence.Words["word1"]; ok {
			fmt.Printf("Found word: %d\n", val)
		}

		// looping over each instance
		for word, occurence := range occurence.Words {

			fmt.Printf("%s: %d\n", word, occurence)

		}
		return occurence, nil
	}

	return nil, nil
}
