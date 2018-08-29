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

import "testing"

func setupOutLogic() *OutgoingLogic {
	return NewOutgoingLogic()
}

func setupHeader(id uint8, mask uint32) Header {
	return Header{SequenceID: NewSequenceID(id), Mask: ReceiveMask{bits: mask}}
}

func TestMultipleReceive(t *testing.T) {
	l := setupOutLogic()
	h := setupHeader(2, 0xffffffff)

	firstErr := l.ReceivedByRemote(h)
	if firstErr != nil {
		t.Error(firstErr)
	}
	if 3 != l.QueueCount() {
		t.Errorf("Count is wrong")
	}

	h2 := setupHeader(4, 0xffffffff)
	secondErr := l.ReceivedByRemote(h2)
	if secondErr != nil {
		t.Error(secondErr)
	}

	if 5 != l.QueueCount() {
		t.Errorf("Wrong count: %v", l.QueueCount())
	}

	for i := 0; i < 5; i++ {
		delivery, err := l.Dequeue()
		if err != nil {
			t.Error(err)
		}
		if !delivery.WasDelivered {
			t.Errorf("Should be delivered")
		}
	}
}

func TestOutgoingSequenceID(t *testing.T) {
	l := setupOutLogic()
	for i := 0; i < ReceiveMaskRange; i++ {
		if !l.CanIncrementOutgoingSequence() {
			t.Errorf("Should be able to increment at this point")
		}
		l.IncreaseOutgoingSequenceID()
	}
	if l.CanIncrementOutgoingSequence() {
		t.Errorf("Should not be able to increment at this point")
	}
	h := setupHeader(2, 0xffffffff)
	l.ReceivedByRemote(h)
	for j := 0; j < 3; j++ {
		if !l.CanIncrementOutgoingSequence() {
			t.Errorf("Should be able to increment at this point")
		}
		l.IncreaseOutgoingSequenceID()
	}
	if l.CanIncrementOutgoingSequence() {
		t.Errorf("Should not be able to increment at this point")
	}
}

func TestSomeDropped(t *testing.T) {
	l := setupOutLogic()
	h := setupHeader(3, 0x2)
	l.ReceivedByRemote(h)
	if 4 != l.QueueCount() {
		t.Errorf("Queue count is wrong %v", l.QueueCount())
	}
	info, infoErr := l.Dequeue()
	if infoErr != nil {
		t.Error(infoErr)
	}
	if info.WasDelivered {
		t.Errorf("Shouldn't be delivered")
	}

	info2, info2Err := l.Dequeue()
	if info2Err != nil {
		t.Error(info2Err)
	}
	if !info2.WasDelivered {
		t.Errorf("Should be delivered")
	}
}
