package katex

import (
	_ "embed"
	"io"

	"github.com/dop251/goja"
)

//go:embed katex.min.js
var katexjs string

func Render(w io.Writer, src []byte, display bool) error {
	var res goja.Value
	vm := goja.New()

	_, err := vm.RunString(katexjs)
	if err != nil {
		return err
	}

	err = vm.Set("expression", string(src))
	if err != nil {
		return nil
	}

	if display {
		res, err = vm.RunString("katex.renderToString(expression, { displayMode: true })")
	} else {
		res, err = vm.RunString("katex.renderToString(expression)")
	}
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, res.Export().(string))
	return err
}
