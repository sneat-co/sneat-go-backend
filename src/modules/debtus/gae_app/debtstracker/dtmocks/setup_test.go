package dtmocks

import (
	"context"
	"testing"
)

func TestSetupMocks(t *testing.T) {
	c := context.Background()
	SetupMocks(c)
}
