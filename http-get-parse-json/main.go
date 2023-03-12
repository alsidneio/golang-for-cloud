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

func main() {

	// JSON to Struct golang
	type Page struct {
		Name string `json:"page"`
	}

	type Words struct {
		Input string   `json:"input"`
		Words []string `json:"words"`
	}

	type Occurrence struct {
		Words map[string]int `json:"words"`
	}

	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: /http-get <url>\n")
		os.Exit(1)
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is in invalid format: %s\n", err)
	}

	response, err := http.Get(args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("invalid output (HTTP code %d): %s\n", response.StatusCode, body)
	}

	var page Page

	json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}
	switch page.Name {
	case "words":
		var words Words

		json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("JSON parsed\nPage: %s\nWords: %v\n", page.Name, strings.Join(words.Words, ", "))

	case "occurrence":
		var occurence Occurrence

		json.Unmarshal(body, &occurence)
		if err != nil {
			log.Fatal(err)
		}

		// Checking if the value exists
		if val, ok := occurence.Words["word1"]; ok {
			fmt.Printf("Found word: %d\n", val)
		}

		// looping over each instance
		for word, occurence := range occurence.Words {

			fmt.Printf("%s: %d\n", word, occurence)

		}

	// fmt.Printf("JSON parsed\nPage: %s\nWords: %v\n", page.Name, )
	default:
		fmt.Println("Page not found")
	}

}
