package domain

import "testing"

func TestRouteModelFields(t *testing.T) {
	r := Route{Code: "ASIA-001", Name: "Asia Loop"}
	if r.Code == "" || r.Name == "" {
		t.Fatal("expected route fields")
	}
}
