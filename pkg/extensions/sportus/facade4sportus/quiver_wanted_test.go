package facade4sportus

import "testing"

func TestDeleteWantedRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     DeleteWantedRequest
		wantErr bool
	}{
		{"valid", DeleteWantedRequest{ID: "w1"}, false},
		{"empty_id", DeleteWantedRequest{ID: ""}, true},
		{"whitespace_id", DeleteWantedRequest{ID: " "}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DeleteWantedRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
