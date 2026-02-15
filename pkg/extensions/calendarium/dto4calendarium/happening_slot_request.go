package dto4calendarium

// HappeningSlotRequest updates slot
type HappeningSlotRequest struct {
	HappeningRequest
	Slot HappeningSlotWithID `json:"slot"`
}

// Validate returns error if not valid
func (v HappeningSlotRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	return nil
}
