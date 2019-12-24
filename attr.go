// Package attributes is a extension for the goldmark
// (http://github.com/yuin/goldmark).
//
// This extension adds support for block attributes in markdowns.
//  {#id .class option="value"}
//  paragraph text with attributes
package attributes

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// attributesBlock parsed block.
type attributesBlock struct {
	ast.BaseBlock
}

// Dump implements Node.Dump.
func (a *attributesBlock) Dump(source []byte, level int) {
	attrs := a.Attributes()
	list := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		var (
			name  = util.BytesToReadOnlyString(attr.Name)
			value = util.BytesToReadOnlyString(util.EscapeHTML(attr.Value.([]byte)))
		)
		list[name] = value
	}
	ast.DumpHelper(a, source, level, list, nil)
}

// kindAttributes is a NodeKind of the Attributes node.
var kindAttributes = ast.NewNodeKind("Attributes")

// Kind implements Node.Kind.
func (a *attributesBlock) Kind() ast.NodeKind {
	return kindAttributes
}

// Attributes defines a markdown block attributes parser, render & extension.
type Attributes struct {
	supportedTypes []ast.NodeKind
}

// DefaultNodeKinds contains a list of the default supported block element
// types.
var DefaultNodeKinds = []ast.NodeKind{
	ast.KindBlockquote, ast.KindHeading, ast.KindList,
	ast.KindParagraph, ast.KindThematicBreak,
	east.KindTable, east.KindDefinitionList,
}

// Extension implement markdown block attributes support.
// Params define a list of supported node types.
// If nil, DefaultNodeKinds are used.
func Extension(nodes ...ast.NodeKind) *Attributes {
	if len(nodes) == 0 {
		nodes = DefaultNodeKinds
	}
	return &Attributes{
		supportedTypes: nodes,
	}
}

// Enable return initialized goldmark.Option with block attributes support.
func Enable(nodes ...ast.NodeKind) goldmark.Option {
	return goldmark.WithExtensions(Extension(nodes...))
}

func (a *Attributes) isSupported(k ast.NodeKind) bool {
	for _, t := range a.supportedTypes {
		if t == k {
			return true
		}
	}
	return false
}

// Transform implement parser.Transformer interface.
func (a *Attributes) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	a.walkAtributes(node)
}

func (a *Attributes) walkAtributes(node ast.Node) {
	for node = node.FirstChild(); node != nil; node = node.NextSibling() {
		if node.Kind() == kindAttributes {
			attrs := node.Attributes() // get attributes
			next := node.NextSibling() // next node
			// remove attributes node
			if p := node.Parent(); p != nil {
				p.RemoveChild(p, node)
			}
			if next == nil {
				break
			}
			node = next
			if node.Type() == ast.TypeBlock &&
				!node.HasBlankPreviousLines() &&
				a.isSupported(node.Kind()) {
				// set attribute to sibling node
				for _, attr := range attrs {
					node.SetAttribute(attr.Name, attr.Value)
				}
			}
		}
		if node.HasChildren() {
			a.walkAtributes(node)
		}
	}
}

// Extend implement goldmark.Extender interface.
func (a *Attributes) Extend(m goldmark.Markdown) {
	if len(a.supportedTypes) == 0 {
		return // nothing to change
	}
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(a, 0)),
		parser.WithASTTransformers(
			util.Prioritized(a, 0),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(a, 0),
		),
	)
}

// RegisterFuncs implement renderer.NodeRenderer interface.
func (a *Attributes) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// not render
	reg.Register(kindAttributes, func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		return ast.WalkSkipChildren, nil
	})
}

// Trigger implement parser.BlockParser interface.
func (a *Attributes) Trigger() []byte {
	return []byte{'{'}
}

var ckAttributes = parser.NewContextKey()

// Open implement parser.BlockParser interface.
func (a *Attributes) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	if attrs, ok := parser.ParseAttributes(reader); ok {
		// add attributes
		var node = &attributesBlock{
			BaseBlock: ast.BaseBlock{},
		}
		for _, attr := range attrs {
			node.SetAttribute(attr.Name, attr.Value)
		}
		// store in context
		list, ok := pc.Get(ckAttributes).([]*attributesBlock)
		if !ok || list == nil {
			list = []*attributesBlock{node}
		} else {
			list = append(list, node)
		}
		pc.Set(ckAttributes, list)
		return node, parser.NoChildren
	}
	return nil, parser.RequireParagraph
}

// Continue implement parser.BlockParser interface.
func (a *Attributes) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

// Close implement parser.BlockParser interface.
func (a *Attributes) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

// CanInterruptParagraph implement parser.BlockParser interface.
func (a *Attributes) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine implement parser.BlockParser interface.
func (a *Attributes) CanAcceptIndentedLine() bool {
	return false
}
