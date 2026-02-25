package docs

import (
	"testing"
)

func TestDocs(t *testing.T) {
	// Swagger docs package
	SwaggerInfo.Title = "Test"
	if SwaggerInfo.Title != "Test" {
		t.Errorf("expected Test")
	}
}
