package domain

import "testing"

func TestUserModelFields(t *testing.T) {
	u := User{Phone: "13800000000"}
	if u.Phone == "" {
		t.Fatal("expected phone")
	}
}
