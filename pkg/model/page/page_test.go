package page

import "testing"

func TestPage(t *testing.T) {
	p := Page[string]{
		Records:    []string{"a", "b"},
		Total:      2,
		PageNum:    1,
		PageSize:   10,
		TotalPages: 1,
	}
	if len(p.Records) != 2 {
		t.Errorf("expected 2 records, got %d", len(p.Records))
	}
	if p.Total != 2 {
		t.Errorf("expected total 2, got %d", p.Total)
	}
}
