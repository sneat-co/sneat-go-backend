package botcmds4splitus

type callbackLink struct {
}

var CallbackLink = callbackLink{}

func (callbackLink) ToGroup(groupID string, isEdit bool) string {
	s := groupCommandCode + "?id=" + groupID
	if isEdit {
		s += "&edit=1"
	}
	return s
}
