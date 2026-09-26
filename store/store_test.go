package store

import (
	"path/filepath"
	"testing"

	"github.com/x1uc/comment_user_analysis/models"
)

func TestNewStoreMigratesAndPreservesData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	first, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}

	info := models.UserPhoneInfo{}
	info.User.IDStr = "123"
	info.User.ScreenName = "test user"
	if err := first.InsertInfo(info); err != nil {
		t.Fatal(err)
	}
	// Simulate a database created before goose tracked migrations.
	if _, err := first.ctx.Exec("DROP TABLE goose_db_version"); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	second, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	if err := second.ctx.QueryRow("SELECT COUNT(*) FROM USER_PHONE_INFO WHERE user_id_str = ?", "123").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("want 1 preserved row, got %d", count)
	}
	if err := second.ctx.QueryRow("SELECT COUNT(*) FROM goose_db_version WHERE version_id = 1 AND is_applied = 1").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("want migration 1 applied once, got %d records", count)
	}
	if err := second.Close(); err != nil {
		t.Fatal(err)
	}

	third, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	defer third.Close()
	if err := third.ctx.QueryRow("SELECT COUNT(*) FROM USER_PHONE_INFO WHERE user_id_str = ?", "123").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("want 1 row after repeated migration, got %d", count)
	}
}
