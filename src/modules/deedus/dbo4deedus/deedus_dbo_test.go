package dbo4deedus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"testing"
)

func TestChallengeDbo(t *testing.T) {
	c := ChallengeDbo{
		Title:  "Challenge 1",
		Status: dbmodels.StatusActive,
		Stars:  5,
	}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
	if c.Title != "Challenge 1" {
		t.Errorf("Title = %v, want Challenge 1", c.Title)
	}
	if c.Stars != 5 {
		t.Errorf("Stars = %v, want 5", c.Stars)
	}

	if err := (ChallengeDbo{}).Validate(); err == nil {
		t.Error("Validate() error = nil, want error for empty title")
	}
}

func TestDeedDbo(t *testing.T) {
	d := DeedDbo{
		Status: dbmodels.StatusActive,
		Starts: 3,
		Details: DeedDetails{
			Text: "Deed details",
		},
	}
	if d.Details.Text != "Deed details" {
		t.Errorf("Details.Text = %v, want Deed details", d.Details.Text)
	}
}
