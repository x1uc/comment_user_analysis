package config

import (
	"log"
	"sync"
)

var (
	globalConfig *Config
	configOnce   sync.Once
)

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig() error {
	var err error
	configOnce.Do(func() {
		globalConfig, err = LoadConfig()
		if err != nil {
			log.Printf("初始化全局配置失败: %v", err)
			return
		}
		log.Println("全局配置初始化成功")
	})
	return err
}

// GetGlobalConfig 获取全局配置单例
func GetGlobalConfig() *Config {
	if globalConfig == nil {
		log.Fatal("全局配置未初始化，请先调用 InitGlobalConfig()")
	}
	return globalConfig
}
