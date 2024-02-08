package misc

import (
	"testing"
)

func TestSlice(t *testing.T) {
	a := make([]struct{ Name string }, 1)
	a[0] = struct{ Name string }{Name: "First"}
	t.Log(a[0])
	a[0].Name = "Second"
	t.Log(a[0])
	b := a[0]
	b.Name = "Third"
	t.Log(a[0])
	t.Log(b)
}
