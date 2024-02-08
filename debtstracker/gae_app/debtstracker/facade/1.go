package facade

import "github.com/strongo/delaying"

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayUpdateUserWithGroups = mustRegisterFunc("UpdateUserWithGroups", delayedUpdateUserWithGroups)
	delayUpdateGroupUsers = mustRegisterFunc("updateGroupUsers", updateGroupUsers)
	delayUpdateContactWithGroups = mustRegisterFunc("UpdateContactWithGroups", delayedUpdateContactWithGroup)
	delayedSetUserReferrer = mustRegisterFunc("setUserReferrer", setUserReferrer)
}

var delayUpdateUserWithGroups delaying.Function
var delayUpdateGroupUsers delaying.Function
var delayUpdateContactWithGroups delaying.Function
var delayedSetUserReferrer delaying.Function
