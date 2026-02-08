package dbo4listus

import "testing"

func TestIsKnownListType(t *testing.T) {
	tests := []struct {
		lt   string
		want bool
	}{
		{ListTypeToDo, true},
		{ListTypeToBuy, true},
		{ListTypeGeneral, true},
		{ListTypeToRead, true},
		{ListTypeToWatch, true},
		{"invalid", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.lt, func(t *testing.T) {
			if got := IsKnownListType(tt.lt); got != tt.want {
				t.Errorf("IsKnownListType(%v) = %v, want %v", tt.lt, got, tt.want)
			}
		})
	}
}

func TestListKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		lk      ListKey
		wantErr bool
	}{
		{"valid_todo", NewListKey(ListTypeToDo, "123"), false},
		{"valid_buy", NewListKey(ListTypeToBuy, "abc"), false},
		{"invalid_format", ListKey("invalid"), true},
		{"invalid_type", ListKey("unknown!123"), true},
		{"empty_id", ListKey("do!"), true},
		{"empty", ListKey(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.lk.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ListKey.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
