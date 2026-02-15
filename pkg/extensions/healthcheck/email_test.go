package healthcheck

import "testing"

func TestTestEmail(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("panic expected")
		}
	}()
	httpGetTestEmail(nil, nil)
}
