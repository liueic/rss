package notifier

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rsswatcher/rsswatcher/internal/parser"
)

const (
	defaultBarkServer = "https://api.day.app"
	notifyTimeout     = 10 * time.Second
)

type BarkNotifier struct {
	deviceKey string
	server    string
	client    *http.Client
}

func NewBark() *BarkNotifier {
	deviceKey := os.Getenv("BARK_DEVICE_KEY")
	server := os.Getenv("BARK_SERVER")
	if server == "" {
		server = defaultBarkServer
	}

	return &BarkNotifier{
		deviceKey: deviceKey,
		server:    server,
		client: &http.Client{
			Timeout: notifyTimeout,
		},
	}
}

func (b *BarkNotifier) Notify(feedName string, items []*parser.Item) error {
	if b.deviceKey == "" {
		return fmt.Errorf("BARK_DEVICE_KEY not set")
	}

	for _, item := range items {
		if err := b.notifyItem(feedName, item); err != nil {
			return err
		}
	}

	return nil
}

func (b *BarkNotifier) notifyItem(feedName string, item *parser.Item) error {
	title := fmt.Sprintf("[%s] %s", feedName, truncate(item.Title, 50))
	body := truncate(item.Description, 100)

	if body == "" {
		body = "New item published"
	}

	opts := map[string]string{
		"group": feedName,
	}

	if item.Link != "" {
		opts["url"] = item.Link
	}

	return b.send(title, body, opts)
}

func (b *BarkNotifier) NotifyAggregate(feedName string, items []*parser.Item) error {
	if b.deviceKey == "" {
		return fmt.Errorf("BARK_DEVICE_KEY not set")
	}

	if len(items) == 0 {
		return nil
	}

	title := fmt.Sprintf("[%s] %d new items", feedName, len(items))

	bodyParts := make([]string, 0, len(items))
	for i, item := range items {
		if i >= 5 {
			bodyParts = append(bodyParts, fmt.Sprintf("... and %d more", len(items)-5))
			break
		}
		bodyParts = append(bodyParts, truncate(item.Title, 60))
	}
	body := strings.Join(bodyParts, "\n")

	opts := map[string]string{
		"group": feedName,
	}

	return b.send(title, body, opts)
}

func (b *BarkNotifier) send(title, body string, opts map[string]string) error {
	u := fmt.Sprintf("%s/%s/%s/%s",
		b.server,
		b.deviceKey,
		url.PathEscape(title),
		url.PathEscape(body),
	)

	q := url.Values{}
	for k, v := range opts {
		q.Set(k, v)
	}
	if queryStr := q.Encode(); queryStr != "" {
		u = u + "?" + queryStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), notifyTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bark API returned status %d", resp.StatusCode)
	}

	return nil
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
