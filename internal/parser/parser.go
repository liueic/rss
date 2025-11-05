package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/mmcdole/gofeed"
)

type Item struct {
	GUID        string
	Link        string
	Title       string
	Description string
	Published   string
	Summary     string // AI生成的总结，可选
}

type Parser struct {
	parser *gofeed.Parser
}

func New() *Parser {
	return &Parser{
		parser: gofeed.NewParser(),
	}
}

func (p *Parser) Parse(data []byte) ([]*Item, error) {
	feed, err := p.parser.ParseString(string(data))
	if err != nil {
		return nil, err
	}

	items := make([]*Item, 0, len(feed.Items))
	for _, feedItem := range feed.Items {
		item := &Item{
			GUID:        feedItem.GUID,
			Link:        feedItem.Link,
			Title:       feedItem.Title,
			Description: cleanDescription(feedItem.Description),
		}

		if feedItem.PublishedParsed != nil {
			item.Published = feedItem.PublishedParsed.Format("2006-01-02 15:04:05")
		} else if feedItem.UpdatedParsed != nil {
			item.Published = feedItem.UpdatedParsed.Format("2006-01-02 15:04:05")
		}

		items = append(items, item)
	}

	return items, nil
}

func cleanDescription(desc string) string {
	desc = strings.TrimSpace(desc)
	runes := []rune(desc)
	if len(runes) <= 200 {
		return desc
	}
	return string(runes[:200]) + "..."
}

func cleanDescriptionBytes(desc string, maxBytes int) string {
	desc = strings.TrimSpace(desc)
	if len(desc) <= maxBytes {
		return desc
	}

	for i := maxBytes; i > 0; i-- {
		if utf8.ValidString(desc[:i]) {
			return desc[:i] + "..."
		}
	}
	return "..."
}
