package signal

import (
	"reflect"

	"github.com/andrebq/organic/cell"
)

type (
	// Matcher is used to match signals sent from cells
	Matcher interface {
		// Match should inspect the input array, match any values,
		// and return the rest so another matcher could be applied
		Match(Tuple) Tuple
	}

	// Tuple is used just to simplify signal processing
	Tuple []interface{}

	matchFn func(Tuple) Tuple
)

// Match applies matchers to the signal and returns true if all items where matched
func Match(s cell.Signal, matchers ...Matcher) bool {
	if len(matchers) == 0 && (len(s.Content()) == 0) {
		return true
	} else if len(matchers) == 0 {
		return false
	}

	var args Tuple
	err := cell.DecodeSignal(&args, s)
	if err != nil {
		// think if decoding errors are worth checking
		return false
	}

	for args.Len() > 0 && len(matchers) > 0 {
		m := matchers[0]
		matchers = matchers[1:]
		nargs := m.Match(args)
		if args.Len() == nargs.Len() {
			return false
		}
		args = nargs
	}

	return len(args) == len(matchers) &&
		len(matchers) == 0
}

// String match string values... quite obvious
func String(sptr *string) Matcher {
	return ptrTo(sptr)
}

// Bytes matches a sequence of bytes
func Bytes(bptr *[]byte) Matcher {
	return ptrTo(bptr)
}

// Receptor matches a string which is cast to a Receptor
//
// Receptors cannot have zero-length
func Receptor(sptr *cell.Receptor) Matcher {
	return matchFn(func(t Tuple) Tuple {
		return t.TakeIf(func(v interface{}) bool {
			b, ok := v.(string)
			if !ok || len(b) == 0 {
				return false
			}
			*sptr = cell.Receptor(b)
			return true
		})
	})
}

func ptrTo(ptr interface{}) Matcher {
	return matchFn(func(t Tuple) Tuple {
		return t.TakeIf(takePtrTo(ptr))
	})
}

// ID matches byte values with a length of 32 bytes
func ID(iptr *cell.ID) Matcher {
	return matchFn(func(t Tuple) Tuple {
		return t.TakeIf(func(v interface{}) bool {
			b, ok := v.([]byte)
			if !ok || len(b) != 32 {
				return false
			}
			*iptr = cell.ID{}
			copy((*iptr)[:], b)
			return true
		})
	})
}

// Drop always matches and simply drops the input
func Drop() Matcher {
	return matchFn(func(_ Tuple) Tuple { return nil })
}

func takePtrTo(ptr interface{}) func(interface{}) bool {
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Type().Kind() != reflect.Ptr {
		return takeFalse
	}

	return func(v interface{}) bool {
		val := reflect.ValueOf(v)
		if !val.Type().AssignableTo(ptrVal.Type().Elem()) {
			return false
		}
		ptrVal.Elem().Set(val)
		return true
	}
}

func takeFalse(_ interface{}) bool { return false }

func (m matchFn) Match(t Tuple) Tuple {
	return m(t)
}
