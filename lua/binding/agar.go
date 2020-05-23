package binding

import (
	"github.com/andrebq/organic/cell"
	lua "github.com/yuin/gopher-lua"
)

type (
	moduleDef map[string]lua.LGFunction
)

func Loader(c cell.Cell) func(*lua.LState) int {
	return func(l *lua.LState) int {
		m := l.SetFuncs(l.NewTable(), module)
		cellUserData := l.NewUserData()
		cellUserData.Value = &c
		l.SetField(m, "cell", cellUserData)
		l.Push(m)
		return 1
	}
}

var module = (moduleDef{})
