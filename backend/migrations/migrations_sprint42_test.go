package migrations

import (
	"os"
	"testing"
)

func TestSprint42MigrationFilesExist(t *testing.T) {
	files := []string{
		"000009_sprint42_cabin_extend.up.sql",
		"000009_sprint42_cabin_extend.down.sql",
		"000010_sprint42_user_passenger_extend.up.sql",
		"000010_sprint42_user_passenger_extend.down.sql",
		"000011_sprint42_order_notify_extend.up.sql",
		"000011_sprint42_order_notify_extend.down.sql",
		"000012_analytics_indexes.up.sql",
		"000012_analytics_indexes.down.sql",
		"000013_shop_info_singleton.up.sql",
		"000013_shop_info_singleton.down.sql",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected %s", f)
		}
	}
}
