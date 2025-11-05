package deduper

import (
	"github.com/rsswatcher/rsswatcher/internal/parser"
	"github.com/rsswatcher/rsswatcher/internal/state"
)

type Deduper struct {
	state *state.State
}

func New(s *state.State) *Deduper {
	return &Deduper{
		state: s,
	}
}

func (d *Deduper) GetNewItems(feedID string, items []*parser.Item, dedupeKey string) []*parser.Item {
	if len(items) == 0 {
		return nil
	}

	lastSeen := d.state.Get(feedID)
	if lastSeen == "" {
		firstItem := items[0]
		d.state.Set(feedID, d.getItemKey(firstItem, dedupeKey))
		return []*parser.Item{firstItem}
	}

	newItems := make([]*parser.Item, 0)
	foundLast := false

	for _, item := range items {
		itemKey := d.getItemKey(item, dedupeKey)
		if itemKey == lastSeen {
			foundLast = true
			break
		}
		newItems = append(newItems, item)
	}

	if len(newItems) > 0 {
		firstItem := newItems[0]
		d.state.Set(feedID, d.getItemKey(firstItem, dedupeKey))
	} else if !foundLast && len(items) > 0 {
		firstItem := items[0]
		d.state.Set(feedID, d.getItemKey(firstItem, dedupeKey))
	}

	return newItems
}

func (d *Deduper) getItemKey(item *parser.Item, dedupeKey string) string {
	switch dedupeKey {
	case "guid":
		if item.GUID != "" {
			return item.GUID
		}
		return item.Link
	case "link":
		return item.Link
	case "title":
		return item.Title
	default:
		if item.GUID != "" {
			return item.GUID
		}
		return item.Link
	}
}
