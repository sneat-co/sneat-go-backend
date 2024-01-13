package models4companius

import "testing"

func TestCompanyBase_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := CompanyBase{}
		if err := v.Validate(); err == nil {
			t.Errorf("should return error")
		}
	})
	t.Run("family", func(t *testing.T) {
		t.Run("should_pass", func(t *testing.T) {
			v := CompanyBase{
				Kind:  "private",
				Type:  "family",
				Title: "",
			}
			if err := v.Validate(); err != nil {
				t.Errorf("should return nil")
			}
		})
		t.Run("error_if_title", func(t *testing.T) {
			v := CompanyBase{
				Kind:  "private",
				Type:  "family",
				Title: "Family",
			}
			if err := v.Validate(); err == nil {
				t.Errorf("should return error")
			}
		})
	})
}
