package dto4contactus

// ContactRequestWithOptionalMessage request
type ContactRequestWithOptionalMessage struct {
	ContactRequest
	Message string `json:"message,omitempty"`
}

// Validate validates request
func (v *ContactRequestWithOptionalMessage) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	return nil
}
