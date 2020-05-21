package cell

import (
	"github.com/fxamacker/cbor"
)

// EncodeSignal takes args and returns a valid Signal which can be sent to another cell
func EncodeSignal(args ...interface{}) (Signal, error) {
	val, err := cbor.Marshal(args, cbor.CTAP2EncOptions())
	if err != nil {
		return Signal{}, err
	}
	return Signal{content: val}, nil
}

// DecodeSignal takes a signal and return the arguments to out
func DecodeSignal(out interface{}, sig Signal) error {
	return cbor.Unmarshal(sig.Content(), out)
}

// RawSignal always returns a valid Signal, albeit empty ones are possible,
// if the input data is not a valid cbor message.
//
// Partial cbor messages are valid but in this case Signal won't contain the entire message
// and error is set as well
func RawSignal(buf []byte) (Signal, error) {
	rest, err := cbor.Valid(buf)
	if err != nil {
		return Signal{}, err
	}
	buf = buf[0 : len(buf)-len(rest)]
	return Signal{content: buf}, nil
}

// RawSignalStr works like RawSignal but expects a string as input
func RawSignalStr(str string) (Signal, error) {
	return RawSignal([]byte(str))
}
