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

func (agent *WeiboAgent) GetHotComments(blogID string, uid string, max_id uint64) ([]models.WeiboComment, uint64, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=0&is_reload=1&id=%s&is_show_bulletin=2&is_mix=0&count=20&uid=%s&fetch_level=0&locale=zh-CN&max_id=%d", blogID, uid, max_id)

	body, err := agent.client.Get(url)
	if err != nil {
		return nil, 0, err
	}

	var response struct {
		Ok    int                   `json:"ok"`
		MaxID uint64                `json:"max_id"`
		Data  []models.WeiboComment `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, 0, err
	}

	return response.Data, response.MaxID, nil
}

func (agent *WeiboAgent) GetNewComments(blogID string, max_id string, uid string) ([]models.WeiboComment, uint64, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=1&is_reload=1&id=%s&is_show_bulletin=2&is_mix=0&count=20&uid=%s&fetch_level=0&locale=zh-CN&max_id=%d", blogID, uid, max_id)

	body, err := agent.client.Get(url)
	if err != nil {
		return nil, 0, err
	}

	var response struct {
		Ok    int                   `json:"ok"`
		MaxID uint64                `json:"max_id"`
		Data  []models.WeiboComment `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, 0, err
	}

	return response.Data, response.MaxID, nil
}
