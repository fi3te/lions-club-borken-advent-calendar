package domain

import (
	"fmt"
	"time"

	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/config"
	"github.com/fi3te/lions-club-borken-advent-calendar/pkg/webpage"
	"golang.org/x/net/html"
)

const SeparatorElement = "hr"

func GetAdventCalendarDoor(cfg *config.Config, timestamp time.Time) (*AdventCalendarDoor, error) {
	doors, err := getAdventCalendar(cfg)
	if err != nil {
		return nil, err
	}
	month := timestamp.Month()
	if month != time.December {
		return nil, fmt.Errorf("expected month 'december', got '%s'", month)
	}
	day := timestamp.Day()
	for _, door := range doors {
		if door.Number == day {
			return door, nil
		}
	}
	return nil, fmt.Errorf("no open advent calendar door found for timestamp '%s'", timestamp)
}

func getAdventCalendar(cfg *config.Config) ([]*AdventCalendarDoor, error) {
	root, err := webpage.GetHtml(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("fetching html document failed: %v", err)
	}
	nodes, err := webpage.FindNodes(root, cfg.ContentCssPath)
	if err != nil {
		return nil, fmt.Errorf("traversing html failed: %v", err)
	}
	if len(nodes) != 1 {
		return nil, fmt.Errorf("expected a single root node, got %d", len(nodes))
	}
	node := nodes[0]
	return readAdventCalendar(node)
}

func readAdventCalendar(n *html.Node) ([]*AdventCalendarDoor, error) {
	nodesForEachDay, err := webpage.SplitDiv(n, SeparatorElement)
	if err != nil {
		return nil, fmt.Errorf("html element could not be split: %v", err)
	}
	adventCalendar := make([]*AdventCalendarDoor, len(nodesForEachDay))
	for i, node := range nodesForEachDay {
		str, err := webpage.RenderNode(node)
		if err != nil {
			return nil, fmt.Errorf("html element could not be rendered: %v", err)
		}
		adventCalendar[i] = &AdventCalendarDoor{
			Number:      i + 1,
			HtmlContent: str,
		}
	}
	return adventCalendar, nil
}
