package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/x1uc/comment_user_analysis/models"
)

type HTTPClient interface {
	Get(url string) ([]byte, error)
}

type WeiboAgent struct {
	client HTTPClient
}

func NewService(c HTTPClient) *WeiboAgent {
	return &WeiboAgent{c}
}

// GetUserInfo fetches user profile information from Weibo.
func (agent *WeiboAgent) GetUserInfo(uid string) (*models.WeiboUser, error) {
	uid = strings.TrimSpace(uid)
	if uid == "" {
		return nil, fmt.Errorf("uid cannot be empty")
	}

	// Weibo AJAX API for profile info
	url := fmt.Sprintf("https://weibo.com/ajax/profile/info?uid=%s", uid)

	body, err := agent.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info for uid %s: %w", uid, err)
	}

	// Weibo API usually returns a wrapper: {"ok": 1, "data": {"user": {...}}}
	var wrapper struct {
		Ok   int `json:"ok"`
		Data struct {
			User models.WeiboUser `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode user info response for uid %s: %w", uid, err)
	}

	if wrapper.Ok != 1 {
		return nil, fmt.Errorf("api error for uid %s: status not ok", uid)
	}

	return &wrapper.Data.User, nil
}

func (*WeiboAgent) GetUserDeteilInfo(uid string) {

}

func (*WeiboAgent) GetUserBlogs(uid string, page int) {

}

func (*WeiboAgent) GetCommentsByPopular(blogID string, max_id string, uid string) {

}

func (*WeiboAgent) GetCommentsByRecently(blogID string, max_id string, uid string) {

}
