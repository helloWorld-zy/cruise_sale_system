package domain

import "testing"

func TestCruiseExtendedFields(t *testing.T) {
	c := Cruise{
		Code:          "HARMONY",
		CrewCount:     2200,
		Length:        362.0,
		Width:         66.0,
		DeckCount:     18,
		RefurbishYear: 2024,
		Status:        2,
	}
	if c.Code == "" {
		t.Fatal("expected code")
	}
	if c.CrewCount == 0 {
		t.Fatal("expected crew_count")
	}
	if c.Length == 0 {
		t.Fatal("expected length")
	}
	if c.Width == 0 {
		t.Fatal("expected width")
	}
	if c.DeckCount == 0 {
		t.Fatal("expected deck_count")
	}
	if c.RefurbishYear == 0 {
		t.Fatal("expected refurbish_year")
	}
	if c.Status != 2 {
		t.Fatal("expected status 2")
	}
}
