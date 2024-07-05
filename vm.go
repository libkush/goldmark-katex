// This code has been borrowed from github.com/bep/gojap, following is the license notice

// MIT License

// Copyright (c) 2022 Bjørn Erik Pedersen

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package katex

import (
	_ "embed"
	"sync"

	"github.com/dop251/goja"
)

//go:embed katex.min.js
var katexjs string

var vmPool = sync.Pool{
	New: func() interface{} {
		vm := goja.New()
		// adding katex lib beforehand
		vm.RunString(katexjs)
		return vm
	},
}

func getVm() *goja.Runtime {
	return vmPool.Get().(*goja.Runtime)
}

func putVm(vm *goja.Runtime) {
	vmPool.Put(vm)
}

// New returns a new Exec.
func New_Exec() *Exec {
	return &Exec{
		pcache: make(map[string]*goja.Program),
	}
}

// Exec is a JavaScript executor that caches compiled programs.
type Exec struct {
	pcache   map[string]*goja.Program
	pcacheMu sync.RWMutex
}

// Arg is a named argument to be passed to RunString.
type Arg struct {
	Name  string
	Value any
}

// RunString compiles and runs the given string s as a JavaScript program.
// Note that the compiled program is cached using the string s as the key.
//
// We reuse VMs across exeuctions. The script in s is compiled in strict mode, but
// other than that it's currerntly the caller's responsibility to ensure that the script is
// not binding any global variables, e.g. by making sure that it's wrapped in a
// function.
func (e *Exec) RunString(s string, args ...Arg) (goja.Value, error) {
	e.pcacheMu.RLock()
	p, ok := e.pcache[s]
	e.pcacheMu.RUnlock()
	if !ok {
		var err error
		p, err = goja.Compile("", s, true)
		if err != nil {
			return nil, err
		}
		e.pcacheMu.Lock()
		e.pcache[s] = p
		e.pcacheMu.Unlock()
	}

	vm := getVm()
	defer func() {
		for _, arg := range args {
			vm.GlobalObject().Delete(arg.Name)
		}
		putVm(vm)
	}()

	for _, arg := range args {
		if err := vm.Set(arg.Name, arg.Value); err != nil {
			return nil, err
		}
	}

	return vm.RunProgram(p)
}

func (e *Exec) MustRunString(s string, args ...Arg) goja.Value {
	v, err := e.RunString(s, args...)
	if err != nil {
		panic(err)
	}
	return v
}
