package pipeline

import (
	"context"
	"log"

	"github.com/x1uc/comment_user_analysis/store"
)

type StorePlugin struct {
	store   *store.Store
	batchID string
}

func NewStorePlugin(db *store.Store, batchID string) *StorePlugin {
	return &StorePlugin{store: db, batchID: batchID}
}

func (p *StorePlugin) Name() string {
	return "store"
}

func (p *StorePlugin) Run(ctx context.Context, mem *Memory) error {
	for _, info := range mem.PhoneInfos {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := p.store.InsertInfo(p.batchID, info); err != nil {
			log.Printf("Failed to insert info for user %s: %v", info.User.ScreenName, err)
		}
	}
	return nil
}
