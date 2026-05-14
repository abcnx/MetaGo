package model

import (
	"testing"

	// 场景1: 引用 pkg/model/page
	"github.com/ACANX/MetaGo/pkg/model/page"

	// 场景2: 引用 pkg/model/rest
	"github.com/ACANX/MetaGo/pkg/model/rest"
)

func TestPageAndRest(t *testing.T) {
	// 使用 page 包 (来自 pkg/model/page)
	p := page.Page[string]{
		Records: []string{"a", "b"},
		Total:   2,
	}

	// 使用 rest 包 (来自 pkg/model/rest)
	resp := rest.Success(p)

	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
}
