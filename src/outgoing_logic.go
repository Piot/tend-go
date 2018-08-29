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

type DeliveryInfo struct {
	SequenceID   SequenceID
	WasDelivered bool
}

type OutgoingLogic struct {
	lastReceivedByRemote SequenceID
	receiveQueue         []DeliveryInfo
	outgoingSequenceID   SequenceID
}

type ReceiveMask struct {
	bits uint32
}

func NewReceiveMask(bits uint32) ReceiveMask {
	return ReceiveMask{bits: bits}
}

func (m ReceiveMask) Bits() uint32 {
	return m.bits
}

const ReceiveMaskRange = 32

type Header struct {
	SequenceID SequenceID
	Mask       ReceiveMask
}

func NewOutgoingLogic() *OutgoingLogic {
	return &OutgoingLogic{lastReceivedByRemote: NewSequenceID(MaxValue), outgoingSequenceID: NewSequenceID(MaxValue)}
}

func (l *OutgoingLogic) pushQueue(sequenceID SequenceID, wasReceived bool) {
	l.receiveQueue = append(l.receiveQueue, DeliveryInfo{WasDelivered: wasReceived, SequenceID: sequenceID})
}

func (l *OutgoingLogic) popQueue() (DeliveryInfo, error) {
	if len(l.receiveQueue) == 0 {
		return DeliveryInfo{}, fmt.Errorf("Queue is empty")
	}
	var popped DeliveryInfo
	popped, l.receiveQueue = l.receiveQueue[0], l.receiveQueue[1:]

	return popped, nil
}

func (l *OutgoingLogic) ReceivedByRemote(header Header) error {
	nextID := header.SequenceID
	if !l.lastReceivedByRemote.IsValidSuccessor(nextID) {
		return fmt.Errorf("unordered packets. Duplicates and old packets should be filtered in other layers")
	}

	distance := l.lastReceivedByRemote.Distance(nextID)
	if distance == 0 {
		return fmt.Errorf("distance should never be zero. It must have been some packet duplication?")
	}

	currentID := l.lastReceivedByRemote
	currentID = currentID.Next()

	bits := NewMutableReceiveMask(header.Mask, distance)
	for i := 0; i < distance; i++ {
		wasReceived, wasReceivedErr := bits.ReadNextBit()
		if wasReceivedErr != nil {
			return wasReceivedErr
		}
		l.pushQueue(currentID, wasReceived)
		currentID = currentID.Next()
	}

	l.lastReceivedByRemote = nextID
	return nil
}

func (l *OutgoingLogic) CanIncrementOutgoingSequence() bool {
	return l.lastReceivedByRemote.Distance(l.outgoingSequenceID) < ReceiveMaskRange
}

func (l *OutgoingLogic) IncreaseOutgoingSequenceID() {
	l.outgoingSequenceID = l.outgoingSequenceID.Next()
}

func (l *OutgoingLogic) OutgoingSequenceID() SequenceID {
	return l.outgoingSequenceID
}

func (l *OutgoingLogic) QueueCount() int {
	return len(l.receiveQueue)
}

func (l *OutgoingLogic) Dequeue() (DeliveryInfo, error) {
	return l.popQueue()
}
