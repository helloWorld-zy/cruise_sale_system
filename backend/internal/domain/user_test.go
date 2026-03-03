package domain

import "testing"

// TestUserModelFields 测试 User 模型字段。
func TestUserModelFields(t *testing.T) {
	u := User{Phone: "13800000000"}
	if u.Phone == "" {
		t.Fatal("expected phone")
	}
}

// TestUserExtendedFields 测试 User 扩展字段。
func TestUserExtendedFields(t *testing.T) {
	u := User{
		Phone:     "13800000000",
		Email:     "test@example.com",
		AlipayUID: "alipay_uid_001",
	}
	if u.Email == "" {
		t.Fatal("expected email")
	}
	if u.AlipayUID == "" {
		t.Fatal("expected alipay_uid")
	}
}
