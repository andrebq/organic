package binding

import (
	"fmt"

	"github.com/andrebq/organic/cell"
	lua "github.com/yuin/gopher-lua"
)

type (
	moduleDef map[string]lua.LGFunction
)

const (
	moduleName   = "organic"
	idTypeName   = moduleName + "_Id"
	cellTypeName = moduleName + "_Cell"
)

func ModuleName() string { return moduleName }

func Loader(c cell.Cell) func(*lua.LState) int {
	return func(l *lua.LState) int {
		m := l.SetFuncs(l.NewTable(), module)
		registerTypes(l, m)

		cellUserData := l.NewUserData()
		cellUserData.Value = &c
		l.SetMetatable(cellUserData, l.GetTypeMetatable(cellTypeName))

		l.SetField(m, "cell", cellUserData)
		l.Push(m)
		return 1
	}
}

func registerTypes(l *lua.LState, module *lua.LTable) {
	registerType(l, cellTypeName, moduleDef{}, cellInstanceMethods)
	l.SetField(module, "id", registerType(l, idTypeName, moduleDef{
		"of":         idOf,
		"__tostring": idToString,
	}, idInstanceMethods))
}

func registerType(l *lua.LState, name string, typeMethods, instanceMethods moduleDef) *lua.LTable {
	mt := l.NewTypeMetatable(name)
	l.SetFuncs(mt, typeMethods)
	l.SetField(mt, "__index", l.SetFuncs(l.NewTable(), instanceMethods))
	return mt
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
	l.Pop(l.GetTop())
	return 0
}

func idToString(l *lua.LState) int {
	st := stateHelper{l}
	id, ok := st.CheckCellID(1)
	if !ok {
		return 0
	}
	l.Pop(l.GetTop())
	l.Push(lua.LString(id.String()))
	return 1
}

func parseId(l *lua.LState) int {
	str := l.CheckString(1)
	l.Pop(l.GetTop())

	id, err := cell.DecodeIDString(str)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	ud := l.NewUserData()
	ud.Value = id
	l.SetMetatable(ud, l.GetTypeMetatable(idTypeName))
	l.Push(ud)
	return 1
}

func idOf(l *lua.LState) int {
	str := l.CheckString(1)
	l.Pop(l.GetTop())
	ud := l.NewUserData()
	ud.Value = cell.ComputeCellID(str)
	l.SetMetatable(ud, l.GetTypeMetatable(idTypeName))
	l.Push(ud)
	return 1
}

var module = moduleDef{
	"parse_id": parseId,
}
var cellInstanceMethods = moduleDef{
	"send": cellSend,
}
var idInstanceMethods = moduleDef{}
