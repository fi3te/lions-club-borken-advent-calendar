package webpage

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func GetHtml(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return node, nil
}
