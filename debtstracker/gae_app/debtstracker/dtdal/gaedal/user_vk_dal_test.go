package gaedal

import (
	"testing"
)

func TestNewVkUserKey(t *testing.T) {
	const vkUserID = 789
	testIntKey(t, vkUserID, NewUserVkKey(vkUserID))
}
