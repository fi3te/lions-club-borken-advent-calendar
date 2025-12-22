package webpage

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

func FindNodes(n *html.Node, path string) ([]*html.Node, error) {
	return findNodes(n, path, "")
}

func findNodes(n *html.Node, path string, pathVisited string) ([]*html.Node, error) {
	if n == nil {
		return nil, errors.New("no node as starting point given")
	}
	if path == "" {
		return nil, errors.New("no path specified")
	}
	pathComponents := strings.Split(path, " ")
	selector := pathComponents[0]
	var foundNodes []*html.Node
	for childNode := range n.ChildNodes() {
		if matches(childNode, selector) {
			foundNodes = append(foundNodes, childNode)
		}
	}

	if len(foundNodes) == 0 || len(pathComponents) == 1 {
		return foundNodes, nil
	} else {
		newPath := strings.Join(pathComponents[1:], " ")
		if pathVisited == "" {
			pathVisited = selector
		} else {
			pathVisited += " " + selector
		}

		var results []*html.Node
		for _, node := range foundNodes {
			if result, err := findNodes(node, newPath, pathVisited); err == nil && result != nil && len(result) > 0 {
				results = append(results, result...)
			}
		}
		return results, nil
	}
}

func matches(n *html.Node, selector string) bool {
	if n.Type != html.ElementNode {
		return false
	} else if strings.Index(selector, "#") != -1 {
		split := strings.Split(selector, "#")
		element := split[0]
		id := split[1]
		return n.Data == element && matchesId(n, id)
	} else {
		return n.Data == selector
	}
}

func matchesId(n *html.Node, id string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "id" && attr.Val == id {
			return true
		}
	}
	return false
}
