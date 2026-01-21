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

	return Store{
		ctx: ctx,
	}
}

func (s Store) InsertInfo(models.UserPhoneInfo) error {
	return nil
}
