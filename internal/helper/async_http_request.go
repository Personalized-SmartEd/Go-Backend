package helper

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 15 * time.Second,
	}
)

type asyncResult struct {
	Body []byte
	Code int
	Err  error
}

func AsyncHTTPRequest(ctx context.Context, method string, url string, body []byte) <-chan asyncResult {
	resultChan := make(chan asyncResult, 1)

	go func() {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
		if err != nil {
			resultChan <- asyncResult{Err: err}
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			resultChan <- asyncResult{Err: err}
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		resultChan <- asyncResult{
			Body: responseBody,
			Code: resp.StatusCode,
			Err:  err,
		}
	}()

	return resultChan
}
