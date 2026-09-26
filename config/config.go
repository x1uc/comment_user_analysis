package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type BlogsInfo struct {
	BlogId        string `toml:"blog_id"`        // 微博文章ID
	CommentAmount int    `toml:"comment_amount"` // 需要拉取多少条数据
	OrderType     string `toml:"order_type"`     // 按照热度排序还是按照时间降序排序(可选值：popular、timeDesc)
}

type BlogCrawlInfos struct {
	BlogsInfos       []BlogsInfo `toml:"blogs_infos"` // 需要拉取的文章信息
	Cookie           string      `toml:"cookie"`
	RateLimit        string      `toml:"rate_limit"`         // 拉取的限流策略
	DefaultOrderType string      `toml:"default_order_type"` // 默认的评论排序方式，在 BlogInfo.OrderType 为空的时候使用
	DBPath           string      `toml:"db_path"`
}

var (
	DefaultDBPath        = "data.db"
	DefaultOrderType     = "popular"
	DefaultRateLimit     = "5s"
	DefaultCommentAmount = 20
)

func Load(path string) (*BlogCrawlInfos, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg BlogCrawlInfos
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	cfg.applyDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}
	return &cfg, nil
}

func (c *BlogCrawlInfos) applyDefaults() {
	if strings.TrimSpace(c.DBPath) == "" {
		c.DBPath = DefaultDBPath
	}
	if strings.TrimSpace(c.RateLimit) == "" {
		c.RateLimit = DefaultRateLimit
	}
	if strings.TrimSpace(c.DefaultOrderType) == "" {
		c.DefaultOrderType = DefaultOrderType
	}
	for i := range c.BlogsInfos {
		if c.BlogsInfos[i].CommentAmount == 0 {
			c.BlogsInfos[i].CommentAmount = DefaultCommentAmount
		}
		if c.BlogsInfos[i].OrderType == "" {
			c.BlogsInfos[i].OrderType = c.DefaultOrderType
		}
	}
}

func (c *BlogCrawlInfos) RateLimitDuration() (time.Duration, error) {
	return time.ParseDuration(c.RateLimit)
}

func (c *BlogCrawlInfos) OrderTypeFor(info BlogsInfo) string {
	if strings.TrimSpace(info.OrderType) != "" {
		return info.OrderType
	}
	return c.DefaultOrderType
}

func (c *BlogCrawlInfos) Validate() error {
	if strings.TrimSpace(c.Cookie) == "" {
		return fmt.Errorf("cookie is required")
	}
	if _, err := time.ParseDuration(c.RateLimit); err != nil {
		return fmt.Errorf("rate_limit %q: %w", c.RateLimit, err)
	}
	if err := validateOrderType(c.DefaultOrderType); err != nil {
		return fmt.Errorf("default_order_type: %w", err)
	}
	if len(c.BlogsInfos) == 0 {
		return fmt.Errorf("blogs_infos is required")
	}
	for i, info := range c.BlogsInfos {
		if strings.TrimSpace(info.BlogId) == "" {
			return fmt.Errorf("blogs_infos[%d].blog_id is required", i)
		}
		if info.CommentAmount <= 0 {
			return fmt.Errorf("blogs_infos[%d].comment_amount must be greater than 0", i)
		}
		if err := validateOrderType(c.OrderTypeFor(info)); err != nil {
			return fmt.Errorf("blogs_infos[%d].order_type: %w", i, err)
		}
	}
	return nil
}

func validateOrderType(orderType string) error {
	switch orderType {
	case "popular", "timeDesc":
		return nil
	default:
		return fmt.Errorf("must be popular or timeDesc, got %q", orderType)
	}
}
