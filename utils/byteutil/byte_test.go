package byteutil

import (
	"container/list"
	"testing"
)

func TestSplitUint8To2Bytes(t *testing.T) {
	t.Log(SplitUint8To2Bytes(26))
}

func TestGenSpecBytes(t *testing.T) {
	t.Log(GenSpecBytes(10))
}

func TestList(t *testing.T) {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	for element, i := l.Front(), 0; element != nil; element, i = element.Next(), i+1 {
		if i == 1 {
			l.InsertBefore(99, element)
		}
	}

	for element := l.Front(); element != nil; element = element.Next() {
		t.Log(element.Value)
	}
}

func TestXOR(t *testing.T) {
	a := 0b0011
	a ^= 0b1100
	t.Logf("%b", ^a)
}
