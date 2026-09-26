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
		return nil, fmt.Errorf("uid 不能为空")
	}

	// Weibo AJAX API for profile info
	url := fmt.Sprintf("https://weibo.com/ajax/profile/info?uid=%s", uid)

	body, err := agent.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取用户 %s 的资料失败: %w", uid, err)
	}

	// Weibo API usually returns a wrapper: {"ok": 1, "data": {"user": {...}}}
	var wrapper struct {
		Ok   int `json:"ok"`
		Data struct {
			User models.WeiboUser `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("解析用户 %s 的资料失败: %w", uid, err)
	}

	if wrapper.Ok != 1 {
		return nil, fmt.Errorf("用户 %s 的接口返回失败", uid)
	}

	return &wrapper.Data.User, nil
}

func (agent *WeiboAgent) GetUserDetailInfo(uid string) (*models.UserDetail, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/profile/detail?uid=%v", uid)
	body, err := agent.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取用户 %s 的详情失败: %w", uid, err)
	}

	var wrapper struct {
		Data models.UserDetail `json:"data"`
		Ok   int               `json:"ok"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("解析用户 %s 的详情失败: %w", uid, err)
	}

	if wrapper.Ok != 1 {
		return nil, fmt.Errorf("用户 %s 的接口返回失败", uid)
	}

	return &wrapper.Data, nil
}

func (agent *WeiboAgent) GetUserBlogs(uid string, page int) ([]models.WeiboBlog, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/mymblog?uid=%s&page=%d&feature=1", uid, page)
	body, err := agent.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取用户 %s 的微博失败: %w", uid, err)
	}

	var wrapper struct {
		Data struct {
			List []models.WeiboBlog `json:"list"`
		} `json:"data"`
		Ok int `json:"ok"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("解析用户 %s 的微博失败: %w", uid, err)
	}

	if wrapper.Ok != 1 {
		return nil, fmt.Errorf("用户 %s 的接口返回失败", uid)
	}
	return wrapper.Data.List, nil
}

func (agent *WeiboAgent) GetHotComments(blogID string, uid string, maxId uint64) ([]models.WeiboComment, uint64, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=0&is_reload=1&id=%s&is_show_bulletin=2&is_mix=0&count=20&uid=%s&fetch_level=0&locale=zh-CN&max_id=%d", blogID, uid, maxId)

	body, err := agent.client.Get(url)
	if err != nil {
		return nil, 0, err
	}

	var wrapper struct {
		Ok    int                   `json:"ok"`
		MaxID uint64                `json:"maxId"`
		Data  []models.WeiboComment `json:"data"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, 0, err
	}

	if wrapper.Ok != 1 {
		return nil, 0, fmt.Errorf("微博 %s 的评论接口返回失败", blogID)
	}

	return wrapper.Data, wrapper.MaxID, nil
}

func (agent *WeiboAgent) GetNewComments(blogID string, uid string, maxId uint64) ([]models.WeiboComment, uint64, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=1&is_reload=1&id=%s&is_show_bulletin=2&is_mix=0&count=20&uid=%s&fetch_level=0&locale=zh-CN&max_id=%d", blogID, uid, maxId)

	body, err := agent.client.Get(url)
	if err != nil {
		return nil, 0, err
	}

	var wrapper struct {
		Ok    int                   `json:"ok"`
		MaxID uint64                `json:"maxId"`
		Data  []models.WeiboComment `json:"data"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, 0, err
	}

	if wrapper.Ok != 1 {
		return nil, 0, fmt.Errorf("微博 %s 的评论接口返回失败", blogID)
	}

	return wrapper.Data, wrapper.MaxID, nil
}
