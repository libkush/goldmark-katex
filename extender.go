package katex

import (
	_ "embed"

	"github.com/bluele/gcache"
	"github.com/dop251/goja"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

//go:embed katex.min.js
var katexjs string

type Extender struct {
}

func (e *Extender) Extend(m goldmark.Markdown) {
	vm := goja.New()
	vm.RunString(katexjs)
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&Parser{}, 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HTMLRenderer{
			cacheInline:  gcache.New(5000).ARC().Build(),
			cacheDisplay: gcache.New(5000).ARC().Build(),
			gojaVM:       vm,
		}, 0),
	))
}
