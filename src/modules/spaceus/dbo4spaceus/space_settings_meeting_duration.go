package dbo4spaceus

// MeetingDurationSettings record
type MeetingDurationSettings struct {
	// All duration are in minutes
	Total     int `json:"total,omitempty" firestore:"total,omitempty"`
	PerMember int `json:"perMember,omitempty" firestore:"perMember,omitempty"`
}

// Validate validates MeetingDurationSettings record
func (v *MeetingDurationSettings) Validate() error {
	return nil
}
