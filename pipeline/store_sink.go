package pipeline

import (
	"context"
	"log/slog"

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
	return "写入数据库"
}

func (p *StorePlugin) Run(ctx context.Context, mem *Memory) error {
	for _, info := range mem.PhoneInfos {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := p.store.InsertInfo(p.batchID, info); err != nil {
			slog.Error("写入用户数据失败", "user", info.User.ScreenName, "error", err)
		}
	}
	return nil
}
