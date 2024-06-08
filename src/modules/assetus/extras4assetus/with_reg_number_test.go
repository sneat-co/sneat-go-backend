package extras4assetus

import (
	"encoding/json"
	"testing"
)

func TestWithOptionalRegNumberFieldJSON(t *testing.T) {
	v := WithOptionalRegNumberField{
		RegNumber: "",
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "{}" {
		t.Fatalf("Unexpected JSON: %s", string(b))
	}
}
