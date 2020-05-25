package cell

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

type (
	// ID contains the id of a given cell
	ID [32]byte

	// Receptor is the name of a receptor in the cell
	Receptor string

	// Signal should contain a valid cbor message
	//
	// content can be extracted by pattern matching
	Signal struct {
		content []byte
	}

	// Membrane defines the interface that a cell can use to interact with the
	// world
	Membrane interface {
		// ExchangeOut sends the given signal to a cell identified with the given CellID
		// and to the provided Receptor
		ExchangeOut(Carrier)

		// Receives a signal one one of the given receptors
		ExchangeIn(*Carrier)
	}

	// Carrier contains the information required cross a Membrane
	Carrier struct {
		From, To ID
		Receptor
		Signal
	}

	// Medium is used by a Cell to build its basic components, most importantly
	// its Membrane
	Medium interface {
		Membrane(ID) Membrane
	}

	// Cell contains a cell definition
	Cell struct {
		id   ID
		name string
		m    Membrane
	}
)

// Grow constructs a new cell and uses the given Membrane.
func Grow(name string, medium Medium) Cell {
	id := ComputeCellID(name)
	c := Cell{
		name: name,
		id:   id,
		m:    medium.Membrane(id),
	}
	return c
}

// ComputeCellID returns the id of a cell with the given name
func ComputeCellID(name string) (id ID) {
	h := sha256.New()
	io.WriteString(h, name)
	h.Sum(id[:0])
	return id
}

// Content returns the signal content
func (s Signal) Content() []byte {
	return s.content
}

// ID returns the ID of this cell
func (c *Cell) ID() ID {
	return c.id
}

// Name of a cell
func (c *Cell) Name() string {
	return c.name
}

// Emit the list of arguments to the given cell/receptor
func (c *Cell) Emit(to ID, recpt Receptor, sig Signal) {
	c.m.ExchangeOut(Carrier{
		From:     c.id,
		To:       to,
		Receptor: recpt,
		Signal:   sig,
	})
}

// Recv receives a signal in one of the given receptors
func (c *Cell) Recv() (ID, Receptor, Signal) {
	var car Carrier
	c.m.ExchangeIn(&car)
	return car.From, car.Receptor, car.Signal
}

// DecodeIDString tries to read a valid ID from a base64 encoded string
//
// The string must match the result of calling id.String
func DecodeIDString(str string) (id ID, err error) {
	var buf []byte
	buf, err = base64.URLEncoding.DecodeString(str)
	if err != nil {
		return
	}
	if len(buf) != len(id[:]) {
		err = errors.New("string MUST represent a 32-byte array")
		return
	}
	copy(id[:], buf)
	return
}

// ID returns the base64 encoding of a given ID
func (i ID) String() string {
	return base64.URLEncoding.EncodeToString(i[:])
}

// Bytes returns a copy of this ID bytes
func (i ID) Bytes() []byte {
	return append([]byte(nil), i[:]...)
}
