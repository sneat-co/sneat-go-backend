package listusbot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListusBotCommands(t *testing.T) {
	assert.Greater(t, len(listusBotCommands), 0)
}
