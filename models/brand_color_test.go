package models

import "testing"

func TestBrandColorsAreUnique(t *testing.T) {
	seen := make(map[string]string)
	for brand, color := range BrandColors {
		if prev, ok := seen[color]; ok {
			t.Fatalf("%s and %s share %s", brand, prev, color)
		}
		seen[color] = brand
	}
}
