package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func makeLoginRequest(client http.Client, requestURL, password string) (string, error) {
	// 1. Generate json object with the loginRequest
	loginRequest := LoginRequest{
		Password: password,
	}
	// Marshall login request to turn it into a json object
	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("marshal error: %s", err)
	}

	// 2. sending an http POST request
	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("http Post error: %s", err)
	}

	// close the response
	defer response.Body.Close()

	// 3. read in the the response for the post request
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	// checking for successful status code
	if response.StatusCode != 200 {
		return "", fmt.Errorf("invalid POST response (HTTP code %d): %s", response.StatusCode, string(resBody))
	}

	// checking for valid json in response
	if !json.Valid(resBody) {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      "No valid JSON returned",
		}
	}

	// 4. Return the Token from the loginResponse if no error in the process above
	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf(" Login Response Unmarshal error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      "Empty Token replied",
		}
	}
	return loginResponse.Token, nil
}
