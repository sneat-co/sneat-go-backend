package dbo4logist

// Milestone is a milestone of an order
type Milestone struct {
	ID   string `json:"id" firestore:"id"`
	Date string `json:"date" firestore:"date"`
}

// WithMilestones is a struct with milestones
type WithMilestones struct {
	Milestones []*Milestone `json:"milestones,omitempty" firestore:"milestones,omitempty"`
}
