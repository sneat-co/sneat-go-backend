package api4listus

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetListRequestParamsFromURL(t *testing.T) {
	u, _ := url.Parse("http://example.com?spaceID=s1&listID=do!123")
	r := &http.Request{URL: u}
	params := getListRequestParamsFromURL(r)

	if string(params.SpaceID) != "s1" {
		t.Errorf("SpaceID = %v, want s1", params.SpaceID)
	}
	if string(params.ListID) != "do!123" {
		t.Errorf("ListID = %v, want do!123", params.ListID)
	}
}
