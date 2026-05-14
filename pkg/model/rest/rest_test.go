package rest

import "testing"

func TestSuccess(t *testing.T) {
	resp := Success("ok")
	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("expected message 'success', got '%s'", resp.Message)
	}
	if resp.Data != "ok" {
		t.Errorf("expected data 'ok', got '%s'", resp.Data)
	}
}

func TestError(t *testing.T) {
	resp := Error[int](500, "internal error")
	if resp.Code != 500 {
		t.Errorf("expected code 500, got %d", resp.Code)
	}
	if resp.Message != "internal error" {
		t.Errorf("expected message 'internal error', got '%s'", resp.Message)
	}
}
