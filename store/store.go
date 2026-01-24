package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/x1uc/comment_user_analysis/models"
)

type Store struct {
	ctx *sql.DB
}

func NewStore(path string) Store {
	log.Printf("Opening database at %s", path)
	ctx, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	create_table_sql := `
	CREATE TABLE IF NOT EXISTS USER_PHONE_INFO (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id_str TEXT,
		screen_name TEXT,
		blog_text_raw TEXT,
		blog_region_name TEXT,
		blog_source TEXT,
		blog_created_at TEXT,
		blog_id_str TEXT,
		blog_mblog_id TEXT,
		user_ip_location TEXT,
		user_created_at TEXT,
		gender TEXT,
		phone_type TEXT,
		phone_brand TEXT
	);`

	_, err = ctx.Exec(create_table_sql)

	if err != nil {
		log.Fatalf("create table error: %v", err)
	}

	return Store{
		ctx: ctx,
	}
}

func (s Store) InsertInfo(user_info models.UserPhoneInfo) error {
	stmt, err := s.ctx.Prepare(`INSERT INTO USER_PHONE_INFO (
		user_id_str,
		screen_name,
		blog_text_raw,
		blog_region_name,
		blog_source,
		blog_created_at,
		blog_id_str,
		blog_mblog_id,
		user_ip_location,
		user_created_at,
		gender,
		phone_type,
		phone_brand
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		user_info.User.IDStr,
		user_info.User.ScreenName,
		user_info.Blog.TextRaw,
		user_info.Blog.RegionName,
		user_info.Blog.Source,
		user_info.Blog.CreatedAt,
		user_info.Blog.IDStr,
		user_info.Blog.MblogID,
		user_info.Detail.IPLocation,
		user_info.Detail.CreatedAt,
		user_info.Detail.Gender,
		user_info.PhoneType,
		user_info.PhoneBrand,
	)
	return err
}
