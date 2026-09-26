package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
cookie = "test-cookie"
rate_limit = "500ms"
default_order_type = "timeDesc"
db_path = "test.db"

[[blogs_infos]]
blog_id = "111"
comment_amount = 20

[[blogs_infos]]
blog_id = "222"
comment_amount = 5
order_type = "popular"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Cookie != "test-cookie" || cfg.DBPath != "test.db" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	d, err := cfg.RateLimitDuration()
	if err != nil {
		t.Fatal(err)
	}
	if d != 500*time.Millisecond {
		t.Fatalf("rate limit = %s", d)
	}
	if got := cfg.OrderTypeFor(cfg.BlogsInfos[0]); got != "timeDesc" {
		t.Fatalf("default order = %s", got)
	}
	if got := cfg.OrderTypeFor(cfg.BlogsInfos[1]); got != "popular" {
		t.Fatalf("blog order = %s", got)
	}
}

func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
cookie = "test-cookie"

[[blogs_infos]]
blog_id = "111"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBPath != DefaultDBPath {
		t.Fatalf("db_path = %s", cfg.DBPath)
	}
	if cfg.RateLimit != DefaultRateLimit {
		t.Fatalf("rate_limit = %s", cfg.RateLimit)
	}
	if cfg.DefaultOrderType != DefaultOrderType {
		t.Fatalf("default_order_type = %s", cfg.DefaultOrderType)
	}
	if cfg.BlogsInfos[0].CommentAmount != DefaultCommentAmount {
		t.Fatalf("comment_amount = %d", cfg.BlogsInfos[0].CommentAmount)
	}
	if got := cfg.OrderTypeFor(cfg.BlogsInfos[0]); got != DefaultOrderType {
		t.Fatalf("order = %s", got)
	}
}

func TestLoadConvertsBase62BlogID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
cookie = "test-cookie"
rate_limit = "1s"
default_order_type = "timeDesc"
db_path = "test.db"

[[blogs_infos]]
blog_base62_id = "Qn4KL6kCN"
comment_amount = 20
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BlogsInfos[0].BlogId != "5254998191509253" {
		t.Fatalf("blog_id = %s", cfg.BlogsInfos[0].BlogId)
	}
}

func TestLoadRejectsInvalidBase62BlogID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
cookie = "test-cookie"
rate_limit = "1s"
default_order_type = "timeDesc"
db_path = "test.db"

[[blogs_infos]]
blog_base62_id = "!!!"
comment_amount = 1
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected invalid blog_base62_id to fail")
	}
}

func TestLoadRejectsInvalidOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
cookie = "test-cookie"
rate_limit = "1s"
default_order_type = "hot"
db_path = "test.db"

[[blogs_infos]]
blog_id = "111"
comment_amount = 1
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected invalid default_order_type to fail")
	}
}
