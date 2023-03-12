package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// JSON to Struct golang
type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	out := []string{}
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurrence))
	}

	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

type Response interface {
	GetResponse() string
}

//========================================Main Function Starts Here========================================//
func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: /http-get <url>\n")
		os.Exit(1)
	}

	res, err := makeRequest(args[1])
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			log.Fatalf("Error: %s (HTTPCode: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
		}
	}

	if res == nil {
		log.Fatalln("No Response")
	}

	fmt.Printf("Response: %s", res.GetResponse())
}

//============================================Application Functions===========================================//

func makeRequest(requestURL string) (Response, error) {

	if _, err := url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("validation error: URL is invalid: %s ", err)
	}

	response, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("http Get error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid output (HTTP code %d): %s\n", response.StatusCode, body)
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
