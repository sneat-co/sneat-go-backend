package email2awsses

import "testing"

func TestNewEmailClient(t *testing.T) {
	if NewEmailClient("", nil) == nil {
		t.Error("NewEmailClient should return value")
	}
}
