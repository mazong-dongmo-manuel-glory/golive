package golive

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DOMElemSelector struct {
	query []string
}

func NewDOMElementSelector() *DOMElemSelector {
	return &DOMElemSelector{
		query: []string{},
	}
}

func (de *DOMElemSelector) setElemen(elemn string) {
	de.query = append(de.query, elemn)
}

func (de *DOMElemSelector) addAttr(key, value string) {
	de.query = append(de.query, "[", key, "=\"", value, "\"]")
}
func (de *DOMElemSelector) toString() string {
	return strings.Join(de.query, "")
}

type DOMSelector struct {
	query []*DOMElemSelector
}

func NewDOMSelector() *DOMSelector {
	return &DOMSelector{
		query: make([]*DOMElemSelector, 0),
	}
}

func (ds *DOMSelector) addChild() *DOMElemSelector {
	de := NewDOMElementSelector()

	ds.query = append(ds.query, de)
	return de
}

func (ds *DOMSelector) addParent() *DOMElemSelector {
	de := NewDOMElementSelector()

	ds.query = append([]*DOMElemSelector{de}, ds.query...)
	return de
}

func (ds *DOMSelector) toString() string {
	e := []string{}

	for _, q := range ds.query {
		e = append(e, q.toString())
	}

	return strings.Join(e, " ")
}

// AttrMapFromNode todo
func AttrMapFromNode(node *html.Node) map[string]string {
	m := map[string]string{}
	for _, attr := range node.Attr {
		m[attr.Key] = attr.Val
	}
	return m
}

// CreateDOMFromString todo
func CreateDOMFromString(data string) (*html.Node, error) {
	reader := bytes.NewReader([]byte(data))

	parent := &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div}

	fragments, err := html.ParseFragmentWithOptions(reader, parent)

	if err != nil {
		return nil, err
	}

	for _, node := range fragments {
		parent.AppendChild(node)
	}

	return parent, nil
}

// RenderNodeToString todo
func RenderNodeToString(e *html.Node) (string, error) {
	var b bytes.Buffer
	err := html.Render(&b, e)

	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// RenderNodesToString todo
func RenderNodesToString(nodes []*html.Node) (string, error) {
	text := ""

	for _, node := range nodes {
		rendered, err := RenderNodeToString(node)

		if err != nil {
			return "", err
		}

		text += rendered
	}

	return text, nil
}

func RenderChildren(parent *html.Node) (string, error) {
	return RenderNodesToString(GetChildrenFromNode(parent))
}

func getClassesSeparated(s string) string {
	return strings.Join(strings.Split(strings.TrimSpace(s), " "), ".")
}

func SelfIndexOfNode(n *html.Node) int {
	ix := 0

	for prev := n.PrevSibling; prev != nil; prev = prev.PrevSibling {
		ix++
	}

	return ix
}

func GetAllChildrenRecursive(n *html.Node) []*html.Node {
	result := make([]*html.Node, 0)

	if n == nil {
		return result
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, c)

		if c != nil {
			result = append(result, GetAllChildrenRecursive(c)...)
		}
	}

	return result
}

// SelectorFromNode
func SelectorFromNode(e *html.Node) (string, error) {

	err := fmt.Errorf("could not provide a valid selector")

	selector := NewDOMSelector()

	if e.Type == html.ElementNode {

		attrs := AttrMapFromNode(e)

		es := selector.addChild()
		es.setElemen("*")

		if attr, ok := attrs["go-live-uid"]; ok {
			es.addAttr("go-live-uid", attr)

			if attr, ok := attrs["key"]; ok {
				es.addAttr("key", attr)
			}
		}
	}

	for parent := e.Parent; parent != nil; parent = parent.Parent {

		attrs := AttrMapFromNode(e)

		es := NewDOMElementSelector()
		es.setElemen("*")

		found := false
		if attr, ok := attrs["go-live-component-id"]; ok {
			es.addAttr("go-live-component-id", attr)
			found = true
		}
		if attr, ok := attrs["go-live-uid"]; ok {
			es.addAttr("go-live-uid", attr)
			found = true
		}

		if attr, ok := attrs["key"]; ok && found {
			es.addAttr("key", attr)
		}

		if !found {
			continue
		}

		return selector.toString(), nil
	}

	return "", err
}

// PathToComponentRoot todo
func PathToComponentRoot(e *html.Node) []int {

	path := make([]int, 0)

	for parent := e; parent != nil; parent = parent.Parent {

		attrs := AttrMapFromNode(parent)

		path = append([]int{SelfIndexOfNode(parent)}, path...)

		if _, ok := attrs["go-live-component-id"]; ok {
			return path
		}
	}

	return path
}
