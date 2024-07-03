package katex

import (
	"io"

	"github.com/dop251/goja"
)

func Render(w io.Writer, src []byte, display bool, vm *goja.Runtime) error {
	var res goja.Value

	err := vm.Set("expression", string(src))
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
