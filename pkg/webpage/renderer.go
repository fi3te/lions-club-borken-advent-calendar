package webpage

import (
	"bytes"

	"golang.org/x/net/html"
)

func RenderNode(n *html.Node) (string, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return "", err
	}
	return buf.String(), nil
}
