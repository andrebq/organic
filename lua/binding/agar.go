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

		mt := l.NewTable()
		l.SetField(mt, "__index", l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
			"send": cellSend,
		}))
		l.SetMetatable(cellUserData, mt)

		l.SetField(m, "cell", cellUserData)
		l.Push(m)
		return 1
	}
}

func cellSend(l *lua.LState) int {
	ud := l.CheckUserData(1)
	cell, ok := ud.Value.(*cell.Cell)
	if !ok {
		l.RaiseError("Expected a organic.Cell object. Are you using cell:send?")
		return 0
	}
	println("Cell ", cell.Name(), " with id ", cell.ID().String())
	return 0
}

var module = (moduleDef{})
