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

const MaxRange = 128
const HalfMaxRange = MaxRange / 2
const MaxValue = 127

var SequenceIDMax = SequenceID{value: MaxValue}

type SequenceID struct {
	value uint8
}

func NewSequenceID(value uint8) SequenceID {
	return SequenceID{value: value}
}

func (s SequenceID) Next() SequenceID {
	nextValue := (s.value + 1) % MaxRange

	return SequenceID{value: nextValue}
}

func IsValid(id uint8) bool {
	return id < MaxRange
}

func (s SequenceID) Distance(other SequenceID) int {
	nextValue := other.value
	idValue := s.value

	if nextValue < idValue {
		nextValue += MaxRange
	}

	diff := nextValue - idValue

	return int(diff)
}

func (s SequenceID) IsValidSuccessor(other SequenceID) bool {
	distance := s.Distance(other)
	return (distance != 0) && (distance < HalfMaxRange)
}

func (s SequenceID) String() string {
	return fmt.Sprintf("[SequenceID %v]", s.value)
}

func (s SequenceID) Value() uint8 {
	return s.value
}
