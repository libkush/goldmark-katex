package katex

import (
	"io"

	"github.com/dop251/goja"
)

func Render(w io.Writer, src []byte, display bool, exec *Exec) error {
	var res goja.Value
	var err error

	if display {
		res, err = exec.RunString("katex.renderToString(expression, { displayMode: true })", Arg{"expression", string(src)})
	} else {
		res, err = exec.RunString("katex.renderToString(expression)", Arg{"expression", string(src)})
	}
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, res.Export().(string))
	return err
}
