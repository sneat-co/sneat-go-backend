package models4brands

type Model struct {
	MakeID   string `json:"makeID" firestore:"makeID"`
	Title    string `json:"title" firestore:"title"`
	FromDate string `json:"from,omitempty" firestore:"from,omitempty"`
	ToDate   string `json:"to,omitempty" firestore:"to,omitempty"`
}
