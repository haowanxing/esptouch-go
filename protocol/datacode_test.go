package protocol

import (
	"testing"
)

func TestNewDataCode(t *testing.T) {
	dc := NewDataCode(26, 0)
	t.Log(dc, dc.GetBytes())
}
