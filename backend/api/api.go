package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

func makeRequest(ctx context.Context, url string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

func handleResponse(ctx context.Context, response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Upstream request failed", "http.response.status_code", response.StatusCode)
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

func MakeAPIRequest[T any](ctx context.Context, url string, header http.Header) (T, error) {

	response, err := makeRequest(ctx, url, header)
	if err != nil {
		return *new(T), err
	}

	body, err := handleResponse(ctx, response)
	if err != nil {
		return *new(T), err
	}

	return parseResponse[T](body)
}
