package email2awsses

type sent struct {
	messageID string
}

func (v sent) MessageID() string {
	return v.messageID
}
