/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package tend

import (
	"fmt"
	"testing"
)

func setupLogic() *IncomingLogic {
	return NewIncomingLogic()
}

func assertEqualInteger(t *testing.T, a uint32, b uint32, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertEqualSequenceID(t *testing.T, a SequenceID, b SequenceID, message string) {
	assertEqualInteger(t, uint32(a.Value()), uint32(b.Value()), message)
}

func assertError(t *testing.T, e error, message string) {
	if e != nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("Should throw error: %v", message)
	}
	t.Fatal(message)
}

func assertSuccess(t *testing.T, e error, message string) {
	if e == nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("Error: %v %v", e, message)
	}
	t.Fatal(message)
}
func TestFirstReceived(t *testing.T) {
	l := setupLogic()
	assertEqualInteger(t, 0, l.ReceivedHeader().Mask.Bits(), "Should be zero before receiving")
	l.ReceivedToUs(NewSequenceID(1))
	assertEqualInteger(t, 1, l.ReceivedHeader().Mask.Bits(), "")
}
func TestDroppedReceived(t *testing.T) {
	l := setupLogic()

	l.ReceivedToUs(NewSequenceID(0))
	h := l.ReceivedHeader()
	if 1 != h.Mask.Bits() {
		t.Errorf("Wrong")
	}

	l.ReceivedToUs(NewSequenceID(2))
	h2 := l.ReceivedHeader()
	if 5 != h2.Mask.Bits() {
		t.Errorf("More wrong")
	}
}

func TestIllegalDistance(t *testing.T) {
	l := setupLogic()
	h := l.ReceivedHeader()

	assertEqualSequenceID(t, SequenceIDMax, h.SequenceID, "")
	assertEqualInteger(t, 0, l.ReceivedHeader().Mask.Bits(), "")
	assertError(t, l.ReceivedToUs(NewSequenceID(32)), "Should be illegal distance")
}

func TestMaximumDistance(t *testing.T) {
	l := setupLogic()
	h := l.ReceivedHeader()

	assertEqualSequenceID(t, SequenceIDMax, h.SequenceID, "")
	assertEqualInteger(t, 0, l.ReceivedHeader().Mask.Bits(), "")
	assertSuccess(t, l.ReceivedToUs(NewSequenceID(0)), "")
	assertSuccess(t, l.ReceivedToUs(NewSequenceID(1)), "")
	assertSuccess(t, l.ReceivedToUs(NewSequenceID(33)), "")
}
