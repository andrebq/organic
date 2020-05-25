package binding

import (
	"fmt"

	"github.com/andrebq/organic/cell"
	lua "github.com/yuin/gopher-lua"
)

type (
	moduleDef map[string]lua.LGFunction
)

func ModuleName() string { return "organic" }

func Loader(c cell.Cell) func(*lua.LState) int {
	return func(l *lua.LState) int {
		m := l.SetFuncs(l.NewTable(), module)
		registerTypes(l, m)

		cellUserData := l.NewUserData()
		cellUserData.Value = &c
		l.SetMetatable(cellUserData, l.GetTypeMetatable(computeTypeName("Cell")))

		l.SetField(m, "cell", cellUserData)
		l.Push(m)
		return 1
	}
}

func registerTypes(l *lua.LState, module *lua.LTable) {
	registerType(l, module, "Cell", moduleDef{}, cellMethods)
}

func registerType(l *lua.LState, module *lua.LTable, name string, typeMethods, instanceMethods moduleDef) {
	mt := l.NewTypeMetatable(computeTypeName(name))
	l.SetFuncs(mt, typeMethods)
	l.SetField(module, name, mt)
	l.SetField(mt, "__index", l.SetFuncs(l.NewTable(), instanceMethods))
}

func computeTypeName(name string) string {
	return fmt.Sprintf("%v_%v", ModuleName(), name)
}

func cellSend(l *lua.LState) int {
	st := stateHelper{l}
	if !st.CheckTopMin(2) {
		return 0
	}
	cell, ok := st.CheckCell(1)
	if !ok {
		return 0
	}

	to, ok := st.CheckCellID(2)
	if !ok {
		return 0
	}
	println("Cell ", cell.Name(), " with id ", cell.ID().String(), " sending to %v", to.String())
	return 0
}

var module = moduleDef{}
var cellMethods = moduleDef{
	"send": cellSend,
}
