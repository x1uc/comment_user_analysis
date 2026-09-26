package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/x1uc/comment_user_analysis/models"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Store struct {
	ctx *sql.DB
}

func NewStore(path string) (*Store, error) {
	log.Printf("Opening database at %s", path)
	ctx, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := ctx.Ping(); err != nil {
		ctx.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	migrationFS, err := fs.Sub(migrations, "migrations")
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("load migrations: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, ctx, migrationFS)
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("initialize migrations: %w", err)
	}
	if _, err := provider.Up(context.Background()); err != nil {
		ctx.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return &Store{ctx: ctx}, nil
}

func (s *Store) Close() error {
	return s.ctx.Close()
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
