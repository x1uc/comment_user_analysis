package models

import "strings"

// BlogResponse 博客列表响应
type BlogResponse struct {
	Data BlogData `json:"data"`
}

// BlogData 博客数据
type BlogData struct {
	List []Blog `json:"list"`
}

// Blog 博客信息
type Blog struct {
	ID        string `json:"idstr"`
	MblogID   string `json:"mblogid"`
	PhoneType string `json:"source"`
	User      User   `json:"user"`
}

// User 用户信息
type User struct {
	ID string `json:"idstr"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	Data  []CommentData `json:"data"`
	MaxID uint64        `json:"max_id"`
}

// CommentData 评论数据
type CommentData struct {
	User CommentUser `json:"user"`
	Text string      `json:"text"`
}

// CommentUser 评论用户
type CommentUser struct {
	ID string `json:"idstr"`
}

// PhoneStatistics 手机统计数据
type PhoneStatistics struct {
	BrandCounts map[string]int `json:"brand_counts"`
	UserCount   int            `json:"user_count"`
}

// StatisticsData 统计数据（用于导出）
type StatisticsData struct {
	PhoneType string `json:"phone_type"`
	Count     int    `json:"count"`
}

// Config 应用配置
type Config struct {
	UID       string `json:"uid"`
	Cookie    string `json:"cookie"`
	Limit     int    `json:"limit"`
	UserAgent string `json:"user_agent"`
}

type AiResponse struct {
	Value int `json:"value"`
}

// PhoneBrandMapping 手机品牌映射
type PhoneBrandMapping map[string]string

// GetBrand 获取手机品牌
func (p PhoneBrandMapping) GetBrand(phoneType string) string {
	brand := strings.TrimSpace(strings.ToLower(phoneType))
	for key, value := range p {
		if strings.Contains(brand, strings.ToLower(key)) {
			return value
		}
	}
	return phoneType // 如果找不到映射，返回原始值
}
