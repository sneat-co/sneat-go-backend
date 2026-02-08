package api4calendarium

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetHappeningRequestParamsFromURL(t *testing.T) {
	u, _ := url.Parse("http://example.com?spaceID=s1&happeningID=h1&happeningType=t1")
	r := &http.Request{URL: u}
	params := getHappeningRequestParamsFromURL(r)

	if string(params.SpaceID) != "s1" {
		t.Errorf("SpaceID = %v, want s1", params.SpaceID)
	}
	if params.HappeningID != "h1" {
		t.Errorf("HappeningID = %v, want h1", params.HappeningID)
	}
	if params.HappeningType != "t1" {
		t.Errorf("HappeningType = %v, want t1", params.HappeningType)
	}
}
