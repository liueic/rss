package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 15 * time.Second
	defaultRetries = 2
	userAgent      = "RSSWatcher/1.0 (+https://github.com/rsswatcher/rsswatcher)"
)

type Fetcher struct {
	client  *http.Client
	retries int
}

func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		retries: defaultRetries,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= f.retries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*attempt) * time.Second
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		data, err := f.fetchOnce(ctx, url)
		if err == nil {
			return data, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("failed after %d retries: %w", f.retries, lastErr)
}

func (f *Fetcher) fetchOnce(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
