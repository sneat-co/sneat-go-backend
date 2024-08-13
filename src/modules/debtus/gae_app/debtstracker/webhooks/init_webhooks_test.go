package webhooks

import "testing"

func TestInitWebhooks(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	InitWebhooks(nil)
}
