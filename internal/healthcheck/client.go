package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrHTTPStatusNotOK = errors.New("HTTP response status is not OK")
)

type Checker interface {
	Check(ctx context.Context, url string) error
}

type checker struct {
	httpClient *http.Client
}

func NewChecker(httpClient *http.Client) Checker {
	return &checker{
		httpClient: httpClient,
	}
}

func (h *checker) Check(ctx context.Context, url string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		return nil
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("%w: %s: %s", ErrHTTPStatusNotOK, response.Status, string(b))
}
