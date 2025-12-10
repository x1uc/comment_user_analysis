package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"single_analysis/internal/utils"
	"time"
)

// Config 应用配置
type Config struct {
	UID        string   `json:"uid"`
	BlogList   []string `json:"blog_list"`
	Cookie     string   `json:"cookie"`
	Limit      int      `json:"limit"`
	Debug      bool     `json:"debug"`
	OutputDir  string   `json:"output_dir"`
	Interval   int      `json:"interval"`
	OutputName string   `json:"output_name"`
	ApiKey     string   `json:"api_key"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	config := &Config{
		Limit:     100,
		Debug:     false,
		OutputDir: "./output",
		Interval:  5,
	}

	// 1. 首先尝试从配置文件加载
	if err := config.loadFromFile(); err != nil {
		fmt.Printf("警告: %v\n", err)
	}

	// 验证配置
	if err := config.validate(); err != nil {
		return nil, err
	}

	// 创建输出目录
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, utils.NewConfigError("创建输出目录失败", err)
	}

	return config, nil
}

// loadFromFile 从配置文件加载配置
func (c *Config) loadFromFile() error {
	// 尝试多个配置文件位置
	configPaths := []string{
		"config.json",
		"config.local.json",
		"./config/config.json",
		"../config.json",
		filepath.Join(os.Getenv("HOME"), ".config", "comment_analyzer", "config.json"),
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return utils.NewConfigError(fmt.Sprintf("读取配置文件 %s 失败", path), err)
			}

			if err := json.Unmarshal(data, c); err != nil {
				return utils.NewConfigError(fmt.Sprintf("解析配置文件 %s 失败", path), err)
			}

			fmt.Printf("从配置文件加载: %s\n", path)
			return nil
		}
	}

	return utils.NewConfigError("未找到配置文件", nil)
}

// validate 验证配置
func (c *Config) validate() error {
	if c.UID == "" {
		return utils.NewConfigError("用户ID不能为空", nil)
	}

	if c.Cookie == "" {
		return utils.NewConfigError("Cookie不能为空", nil)
	}

	if c.Limit <= 0 {
		c.Limit = 100
	}

	return nil
}

// Save 保存配置到文件
func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return utils.NewConfigError("序列化配置失败", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadFromFile 从文件加载配置
func LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, utils.NewConfigError("读取配置文件失败", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, utils.NewConfigError("解析配置文件失败", err)
	}

	return &config, nil
}

// Print 打印配置信息
func (c *Config) Print() {
	fmt.Printf("配置信息:\n")
	fmt.Printf("  用户ID: %s\n", c.UID)

	// 只显示Cookie的前8个字符，保护隐私
	cookieDisplay := c.Cookie
	if len(cookieDisplay) > 8 {
		cookieDisplay = cookieDisplay[:8] + "..."
	}
	fmt.Printf("  Cookie: %s\n", cookieDisplay)

	fmt.Printf("  统计限制: %d\n", c.Limit)
	fmt.Printf("  调试模式: %t\n", c.Debug)
	fmt.Printf("  输出目录: %s\n", c.OutputDir)
	fmt.Printf("  间隔时间: %d\n", c.Interval)
	fmt.Printf("  开始时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()
}
