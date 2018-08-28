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

import "fmt"

type IncomingLogic struct {
	lastReceivedToUs SequenceID
	receiveMask      uint32
}

func NewIncomingLogic() *IncomingLogic {
	return &IncomingLogic{lastReceivedToUs: NewSequenceID(MaxValue), receiveMask: 0}
}

func (l *IncomingLogic) ReceivedToUs(nextID SequenceID) error {
	if !l.lastReceivedToUs.IsValidSuccessor(nextID) {
		return fmt.Errorf("Unordered packets. Duplicates and old packets should be filtered in other layers")
	}

	distance := l.lastReceivedToUs.Distance(nextID)
	if distance == 0 {
		return fmt.Errorf("Distance should not be zero. Packet duplicates?")
	}

	if distance > ReceiveMaskRange {
		return fmt.Errorf("too big gap in sequence")
	}

	for i := 0; i < distance-1; i++ {
		l.appendReceived(false)
	}
	l.appendReceived(true)
	l.lastReceivedToUs = nextID
	return nil
}

func (l *IncomingLogic) appendReceived(wasReceived bool) {
	l.receiveMask <<= 1
	bits := uint32(0x0)
	if wasReceived {
		bits = uint32(0x01)
	}
	l.receiveMask |= bits
}

func (l *IncomingLogic) ReceivedHeader() Header {
	return Header{SequenceID: l.lastReceivedToUs, Mask: ReceiveMask{bits: l.receiveMask}}
}
