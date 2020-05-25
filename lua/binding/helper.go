package binding

import (
	"fmt"

	"github.com/andrebq/organic/cell"
	lua "github.com/yuin/gopher-lua"
)

type (
	stateHelper struct{ *lua.LState }
)

func (l stateHelper) TypeError(n int, expected string) {
	l.RaiseError(fmt.Sprintf("Expecting argument %v to be %v UserData", n, expected))
}

func (l stateHelper) CheckCell(n int) (*cell.Cell, bool) {
	ud := l.CheckUserData(n)
	cell, ok := ud.Value.(*cell.Cell)
	if !ok {
		l.TypeError(n, "organic.Cell")
	}
	return cell, ok
}

func (l stateHelper) CheckCellID(n int) (cell.ID, bool) {
	ud := l.CheckUserData(n)
	id, ok := ud.Value.(cell.ID)
	if !ok {
		l.TypeError(n, "organic.CellID")
	}
	return id, ok
}

func (l stateHelper) CheckTopMin(n int) bool {
	t := l.GetTop()
	if t < n {
		l.RaiseError(fmt.Sprintf("Expecting %v arguments", n))
		return false
	}
	return true
}
