package dto4calendarium

// DeleteHappeningSlotRequest updates slot
type DeleteHappeningSlotRequest struct {
	HappeningSlotRefRequest
}

// Validate returns error if not valid
func (v DeleteHappeningSlotRequest) Validate() error {
	if err := v.HappeningSlotRefRequest.Validate(); err != nil {
		return err
	}
	return nil
}
