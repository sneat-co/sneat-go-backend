package dbo4logist

import (
	"testing"

	"github.com/sneat-co/contactus-ext/backend/contactusmodels/briefs4contactus"
	"github.com/stretchr/testify/assert"
)

func orderContact(id, parentID string) *OrderContact {
	return &OrderContact{ID: id, Type: briefs4contactus.ContactTypeCompany, Title: "T-" + id, ParentID: parentID, CountryID: "US"}
}

func TestOrderContact_String(t *testing.T) {
	s := OrderContact{ID: "c1", Type: briefs4contactus.ContactTypeCompany, ParentID: "p1", Title: "Title1"}.String()
	assert.Contains(t, s, "c1")
	assert.Contains(t, s, "p1")
	assert.Contains(t, s, "Title1")
}

func TestOrderContact_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderContact
		wantErr bool
	}{
		{"valid", OrderContact{ID: "c1", Type: briefs4contactus.ContactTypeCompany, Title: "T1", CountryID: "US"}, false},
		{"missing_id", OrderContact{Type: briefs4contactus.ContactTypeCompany, Title: "T1", CountryID: "US"}, true},
		{"bad_type", OrderContact{ID: "c1", Type: "bad", Title: "T1", CountryID: "US"}, true},
		{"parent_equals_id", OrderContact{ID: "c1", Type: briefs4contactus.ContactTypeCompany, Title: "T1", CountryID: "US", ParentID: "c1"}, true},
		{"missing_title", OrderContact{ID: "c1", Type: briefs4contactus.ContactTypeCompany, CountryID: "US"}, true},
		{"missing_country", OrderContact{ID: "c1", Type: briefs4contactus.ContactTypeCompany, Title: "T1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithOrderContacts_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Error(t, WithOrderContacts{}.Validate())
	})
	t.Run("valid", func(t *testing.T) {
		v := WithOrderContacts{Contacts: []*OrderContact{orderContact("c1", ""), orderContact("c2", "c1")}}
		assert.NoError(t, v.Validate())
	})
	t.Run("invalid_child", func(t *testing.T) {
		v := WithOrderContacts{Contacts: []*OrderContact{{ID: "c1"}}}
		assert.Error(t, v.Validate())
	})
	t.Run("missing_parent", func(t *testing.T) {
		v := WithOrderContacts{Contacts: []*OrderContact{orderContact("c1", "absent")}}
		assert.Error(t, v.Validate())
	})
}

func TestWithOrderContacts_Updates(t *testing.T) {
	u := WithOrderContacts{Contacts: []*OrderContact{orderContact("c1", "")}}.Updates()
	assert.Len(t, u, 1)
	assert.Equal(t, "contacts", u[0].FieldName())
}

func TestWithOrderContacts_GetContactByID(t *testing.T) {
	c1 := orderContact("c1", "")
	v := WithOrderContacts{Contacts: []*OrderContact{c1}}
	i, c := v.GetContactByID("c1")
	assert.Equal(t, 0, i)
	assert.Equal(t, c1, c)
	i, c = v.GetContactByID("missing")
	assert.Equal(t, -1, i)
	assert.Nil(t, c)
}

func TestWithOrderContacts_MustGetContactByID(t *testing.T) {
	c1 := orderContact("c1", "")
	v := WithOrderContacts{Contacts: []*OrderContact{c1}}
	assert.Equal(t, c1, v.MustGetContactByID("c1"))
	assert.Panics(t, func() { v.MustGetContactByID("missing") })
	assert.Panics(t, func() { v.MustGetContactByID("  ") })
}

func TestWithOrderContacts_GetContactByParentID(t *testing.T) {
	c2 := orderContact("c2", "c1")
	v := WithOrderContacts{Contacts: []*OrderContact{orderContact("c1", ""), c2}}
	i, c := v.GetContactByParentID("c1")
	assert.Equal(t, 1, i)
	assert.Equal(t, c2, c)
	i, c = v.GetContactByParentID("missing")
	assert.Equal(t, -1, i)
	assert.Nil(t, c)
}
