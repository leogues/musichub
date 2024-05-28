package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func makeRequest(url string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func handleResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		log.Print("Erro na solicitação: " + string(body))
		return nil, errors.New("Erro na solicitação: " + string(body))
	}

	return body, nil
}

func parseResponse[T any](body []byte) (T, error) {
	var result T
	err := json.Unmarshal(body, &result)
	if err != nil {
		return *new(T), err
	}

	return result, nil
}

func MakeAPIRequest[T any](url string, header http.Header) (T, error) {

	response, err := makeRequest(url, header)
	if err != nil {
		return *new(T), err
	}

	body, err := handleResponse(response)
	if err != nil {
		return *new(T), err
	}

	return parseResponse[T](body)
}
