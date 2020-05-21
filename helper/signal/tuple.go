package signal

// First returns the first item and the tail
func (t Tuple) First() (interface{}, Tuple) {
	return t[0], t[1:]
}

// TakeIf takes the first item and applies fn to it
// if the function returns true then it returns the rest
// otherwise returns itself
func (t Tuple) TakeIf(fn func(interface{}) bool) Tuple {
	first, tail := t.First()
	if fn(first) {
		return tail
	}
	return t
}

// Len returns the number of items in the tuple
func (t Tuple) Len() int { return len(t) }
