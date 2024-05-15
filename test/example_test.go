package test

import "testing"

func TestExample(t *testing.T) {
	result := "A"
	if result != "A" {
		t.Fatalf("Seharusnya A")
	}
}
