package services

import (
	"testing"
	"time"
)

func TestPhone(t *testing.T) {
	for _, v := range []string{"", "+371 29 580 232", "+1 (202) 555.0123"} {
		if !ValidPhone(v) {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"------", "123456", "1234567890123456", "abc1234567"} {
		if ValidPhone(v) {
			t.Fatal(v)
		}
	}
}
func TestRigaDateBoundary(t *testing.T) {
	now := time.Date(2026, 9, 14, 21, 5, 0, 0, time.UTC)
	if ValidDate("2026-09-14", now) || !ValidDate("2026-09-15", now) {
		t.Fatal("Riga boundary")
	}
	if ValidDate("2026-02-30", now) {
		t.Fatal("invalid date")
	}
}
