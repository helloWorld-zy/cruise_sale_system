package domain

import "testing"

// TestUserModelFields 测试 User 模型字段。
func TestUserModelFields(t *testing.T) {
	u := User{Phone: "13800000000"}
	if u.Phone == "" {
		t.Fatal("expected phone")
	}
}
