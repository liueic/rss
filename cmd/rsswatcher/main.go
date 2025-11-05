package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/rsswatcher/rsswatcher/internal/config"
	"github.com/rsswatcher/rsswatcher/internal/deduper"
	"github.com/rsswatcher/rsswatcher/internal/env"
	"github.com/rsswatcher/rsswatcher/internal/fetcher"
	"github.com/rsswatcher/rsswatcher/internal/notifier"
	"github.com/rsswatcher/rsswatcher/internal/parser"
	"github.com/rsswatcher/rsswatcher/internal/state"
	"github.com/rsswatcher/rsswatcher/internal/summarizer"
)

const maxConcurrent = 8

func main() {
	// 加载 .env 文件（如果存在）
	// 这允许本地开发时使用 .env 文件，而不影响生产环境
	if err := env.LoadEnvDefault(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	configPath := flag.String("config", "feeds.yaml", "Path to feeds configuration file")
	statePath := flag.String("state", "state/last_states.json", "Path to state file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Feeds) == 0 {
		log.Println("No feeds configured")
		return
	}

	// Load state
	s, err := state.Load(*statePath)
	if err != nil {
		log.Fatalf("Failed to load state: %v", err)
	}

	// Initialize components
	fetcher := fetcher.New()
	parser := parser.New()
	deduper := deduper.New(s)
	notifier := notifier.NewBark()
	summarizer := summarizer.New()

	// Process feeds concurrently with semaphore
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, feed := range cfg.Feeds {
		// Set defaults
		if feed.DedupeKey == "" {
			feed.DedupeKey = "guid"
		}
		if feed.Notify {
			feed.Notify = true
		}

		wg.Add(1)
		go func(f config.Feed) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			processFeed(context.Background(), f, fetcher, parser, deduper, notifier, summarizer)
		}(feed)
	}

	wg.Wait()

	// Save state
	if err := s.Save(*statePath); err != nil {
		log.Printf("Failed to save state: %v", err)
	} else {
		log.Printf("State saved to %s", *statePath)
	}
}

func processFeed(ctx context.Context, feed config.Feed, fetcher *fetcher.Fetcher, parser *parser.Parser, deduper *deduper.Deduper, notifier *notifier.BarkNotifier, summarizer *summarizer.Summarizer) {
	log.Printf("Processing feed: %s (%s)", feed.Name, feed.ID)

	// Fetch feed
	data, err := fetcher.Fetch(ctx, feed.URL)
	if err != nil {
		log.Printf("Failed to fetch %s: %v", feed.Name, err)
		return
	}

	// Parse feed
	items, err := parser.Parse(data)
	if err != nil {
		log.Printf("Failed to parse %s: %v", feed.Name, err)
		return
	}

	if len(items) == 0 {
		log.Printf("No items found in %s", feed.Name)
		return
	}

	// Deduplicate
	newItems := deduper.GetNewItems(feed.ID, items, feed.DedupeKey)
	if len(newItems) == 0 {
		log.Printf("No new items in %s", feed.Name)
		return
	}

	log.Printf("Found %d new items in %s", len(newItems), feed.Name)

	// Generate summaries if enabled
	if summarizer.IsEnabled() {
		log.Printf("Generating summaries for %s...", feed.Name)
		for _, item := range newItems {
			summary, err := summarizer.Summarize(ctx, item.Title, item.Description)
			if err != nil {
				log.Printf("Failed to generate summary for %s: %v, using original description", item.Title, err)
				// 总结失败时使用原始描述
				item.Summary = ""
			} else {
				item.Summary = summary
				log.Printf("Generated summary for: %s", item.Title)
			}
		}
	}

	// Send notifications
	if !feed.Notify {
		log.Printf("Notifications disabled for %s", feed.Name)
		return
	}

	if feed.Aggregate {
		if err := notifier.NotifyAggregate(feed.Name, newItems); err != nil {
			log.Printf("Failed to send aggregate notification for %s: %v", feed.Name, err)
		} else {
			log.Printf("Sent aggregate notification for %s (%d items)", feed.Name, len(newItems))
		}
	} else {
		if err := notifier.Notify(feed.Name, newItems); err != nil {
			log.Printf("Failed to send notifications for %s: %v", feed.Name, err)
		} else {
			log.Printf("Sent %d notifications for %s", len(newItems), feed.Name)
		}
	}
}
