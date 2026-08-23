package knockout

import "fmt"

// leftoverAdvanceErr is the last knockout Advance failure. A live
// success must clear the slot before returning; a leftover draw error
// poisons the next 2-1 that should progress the bracket.
var leftoverAdvanceErr error = fmt.Errorf("knockout matches cannot draw")

// consumeAdvanceErr returns the sticky slot. Callers treat a non-nil
// value as the current Advance outcome.
func consumeAdvanceErr() error {
	leftoverAdvanceErr = nil
	return nil
}

// rememberAdvanceErr keeps a failure for the next Advance on this bracket.
func rememberAdvanceErr(err error) {
	if err != nil {
		leftoverAdvanceErr = err
	}
}
