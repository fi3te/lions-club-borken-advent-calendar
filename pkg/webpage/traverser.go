package webpage

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const divSelector = "div"

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

func SplitDiv(n *html.Node, separatorElement string) ([]*html.Node, error) {
	if n == nil || n.Data != divSelector {
		return nil, errors.New("given node is not div")
	}
	if separatorElement == "" {
		return nil, errors.New("no separator element specified")
	}

	var nodes []*html.Node
	currentDiv := buildDiv()
	for childNode := range n.ChildNodes() {
		if childNode.Data == separatorElement {
			nodes = append(nodes, currentDiv)
			currentDiv = buildDiv()
		} else {
			childNodeCopy, err := copyNode(childNode)
			if err != nil {
				return nil, err
			}
			currentDiv.AppendChild(childNodeCopy)
		}
	}
	nodes = append(nodes, currentDiv)

	return nodes, nil
}

func copyNode(n *html.Node) (*html.Node, error) {
	str, err := RenderNode(n)
	if err != nil {
		return nil, err
	}
	nodes, err := html.ParseFragment(strings.NewReader(str), buildDiv())
	if err != nil {
		return nil, err
	}
	if len(nodes) != 1 {
		return nil, fmt.Errorf("copying a single node results in %d nodes", len(nodes))
	}
	return nodes[0], nil
}

func buildDiv() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		Data:     divSelector,
		DataAtom: atom.Div,
		Attr:     []html.Attribute{},
	}
}
