package model

import (
	"testing"

	// 场景3: 两个同名包 page，使用别名区分
	modelpage "github.com/ACANX/MetaGo/pkg/model/page"
	basepage "github.com/ACANX/MetaGo/pkg/base/page"
)

func TestSameNamePackages(t *testing.T) {
	// 使用别名 basepage 引用 pkg/base/page
	paginator := basepage.Paginator{
		CurrentPage: 2,
		PerPage:     10,
	}
	if paginator.Offset() != 10 {
		t.Errorf("expected offset 10, got %d", paginator.Offset())
	}

	// 使用别名 modelpage 引用 pkg/model/page
	p := modelpage.Page[int]{
		Records: []int{1, 2, 3},
		Total:   3,
	}
	if p.Total != 3 {
		t.Errorf("expected total 3, got %d", p.Total)
	}
}
