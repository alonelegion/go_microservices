package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "espresso",
		Price: 1.00,
		SKU:   "abc-acb-bac",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
