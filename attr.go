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
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Block are parsed attributes block.
type Block struct {
	ast.BaseBlock
}

// New return new attributes block.
func New() *Block {
	return &Block{
		BaseBlock: ast.BaseBlock{},
	}
}

// Dump implements Node.Dump.
func (a *Block) Dump(source []byte, level int) {
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

// KindAttributes is a NodeKind of the attributes block node.
var KindAttributes = ast.NewNodeKind("Attributes")

// Kind implements Node.Kind.
func (a *Block) Kind() ast.NodeKind {
	return KindAttributes
}

type attrParser struct{}

var defaultParser = new(attrParser)

// NewParser return new attributes block parser.
func NewParser() parser.BlockParser {
	return defaultParser
}

// Trigger implement parser.BlockParser interface.
func (a *attrParser) Trigger() []byte {
	return []byte{'{'}
}

// Open implement parser.BlockParser interface.
func (a *attrParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	if attrs, ok := parser.ParseAttributes(reader); ok {
		// add attributes
		var node = New()
		for _, attr := range attrs {
			node.SetAttribute(attr.Name, attr.Value)
		}
		return node, parser.NoChildren
	}
	return nil, parser.RequireParagraph
}

// Continue implement parser.BlockParser interface.
func (a *attrParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

// Close implement parser.BlockParser interface.
func (a *attrParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

// CanInterruptParagraph implement parser.BlockParser interface.
func (a *attrParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine implement parser.BlockParser interface.
func (a *attrParser) CanAcceptIndentedLine() bool {
	return false
}

type transformer struct{}

var defaultTransformer = new(transformer)

// NewTransformer return new AST Transformer for attributes block.
func NewTransformer() parser.ASTTransformer {
	return defaultTransformer
}

// Transform implement parser.Transformer interface.
func (a *transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	// collect all attributes block
	var attributes = make([]ast.Node, 0, 1000)
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == KindAttributes {
			attributes = append(attributes, node)
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	// set attributes to next block sibling
	for _, attr := range attributes {
		if next := attr.NextSibling(); next != nil &&
			next.Type() == ast.TypeBlock &&
			!next.HasBlankPreviousLines() {
			// set attribute to sibling node
			for _, attr := range attr.Attributes() {
				if _, exist := next.Attribute(attr.Name); !exist {
					next.SetAttribute(attr.Name, attr.Value)
				}
			}
		}
		// remove attributes node
		attr.Parent().RemoveChild(attr.Parent(), attr)
	}
}

type attrRender struct{}

var defaultRenderer = new(attrRender)

// NewRenderer return new attributes block renderer.
func NewRenderer() renderer.NodeRenderer {
	return defaultRenderer
}

// RegisterFuncs implement renderer.NodeRenderer interface.
func (a *attrRender) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// not render
	reg.Register(KindAttributes,
		func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
			return ast.WalkSkipChildren, nil
		})
}

// extension defines a goldmark.Extender for markdown block attributes.
type extension struct{}

// Extend implement goldmark.Extender interface.
func (a *extension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(defaultParser, 100)),
		parser.WithASTTransformers(
			util.Prioritized(defaultTransformer, 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(defaultRenderer, 100),
		),
	)
}

// Extension is a goldmark.Extender with markdown block attributes support.
var Extension goldmark.Extender = new(extension)

// Enable is a goldmark.Option with block attributes support.
var Enable = goldmark.WithExtensions(Extension)
