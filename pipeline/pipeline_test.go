package pipeline

import (
	"context"
	"errors"
	"testing"

	"github.com/x1uc/comment_user_analysis/models"
)

type fakePlugin struct {
	name  string
	calls *[]string
	run   func(mem *Memory) error
	ran   bool
}

func (p *fakePlugin) Name() string { return p.name }

func (p *fakePlugin) Run(ctx context.Context, mem *Memory) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.ran = true
	*p.calls = append(*p.calls, p.name)
	if p.run != nil {
		return p.run(mem)
	}
	return nil
}

func TestRunOrderAndSharedMemory(t *testing.T) {
	var calls []string
	first := &fakePlugin{
		name:  "comment_user",
		calls: &calls,
		run: func(mem *Memory) error {
			mem.Users = append(mem.Users, models.WeiboUser{IDStr: "1"})
			mem.Comments = append(mem.Comments, models.WeiboComment{IDStr: "c1"})
			return nil
		},
	}
	second := &fakePlugin{
		name:  "phone_info",
		calls: &calls,
		run: func(mem *Memory) error {
			if len(mem.Users) != 1 || mem.Users[0].IDStr != "1" {
				t.Fatalf("users = %+v", mem.Users)
			}
			mem.PhoneInfos = append(mem.PhoneInfos, models.UserPhoneInfo{
				User: mem.Users[0],
			})
			return nil
		},
	}
	third := &fakePlugin{
		name:  "store",
		calls: &calls,
		run: func(mem *Memory) error {
			if len(mem.PhoneInfos) != 1 || mem.PhoneInfos[0].User.IDStr != "1" {
				t.Fatalf("phone infos = %+v", mem.PhoneInfos)
			}
			return nil
		},
	}

	err := Run(context.Background(), []Plugin{first, second, third})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(calls), 3; got != want {
		t.Fatalf("calls = %v", calls)
	}
	if calls[0] != "comment_user" || calls[1] != "phone_info" || calls[2] != "store" {
		t.Fatalf("order = %v", calls)
	}
}

func TestRunStopsAfterError(t *testing.T) {
	var calls []string
	first := &fakePlugin{name: "comment_user", calls: &calls}
	second := &fakePlugin{
		name:  "phone_info",
		calls: &calls,
		run: func(mem *Memory) error {
			return errors.New("boom")
		},
	}
	third := &fakePlugin{name: "store", calls: &calls}

	err := Run(context.Background(), []Plugin{first, second, third})
	if err == nil {
		t.Fatal("expected error")
	}
	if third.ran {
		t.Fatal("store plugin ran after an earlier error")
	}
	if len(calls) != 2 || calls[0] != "comment_user" || calls[1] != "phone_info" {
		t.Fatalf("calls = %v", calls)
	}
}
