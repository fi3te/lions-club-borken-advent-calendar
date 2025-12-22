package webpage

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const divSelector = "div"

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

func AddToNewDiv(n1 *html.Node, n2 *html.Node) *html.Node {
	div := buildDiv()
	div.AppendChild(n1)
	div.AppendChild(n2)
	return div
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
