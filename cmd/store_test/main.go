package main

import (
	"github.com/x1uc/comment_user_analysis/models"
	"github.com/x1uc/comment_user_analysis/store"
)


func main() {
	store := store.NewStore("test.db")
	store.InsertInfo(models.UserPhoneInfo{})
}
