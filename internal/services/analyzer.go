package services

import (
	"fmt"
	"os"
	"path/filepath"
	"single_analysis/config"
	"single_analysis/internal/llm"
	"single_analysis/internal/models"
	"sort"
	"strings"
	"sync"
	"time"
)

// AnalyzerService 分析服务
type AnalyzerService struct {
	weiboService   *WeiboService
	statistics     *models.PhoneStatistics
	statisticsLLM  map[string]models.StatisticData
	processedUsers map[string]bool // 存储已处理过的用户ID，避免重复处理
	outputDir      string          // 输出目录
	statsFile      *os.File        // 实时统计数据文件
	statsFileLLM   *os.File
	mutex          sync.RWMutex
	interval       int
	LLMAnalysis    llm.LLMClient
}

// NewAnalyzerService 创建分析服务
func NewAnalyzerService(weiboService *WeiboService, cfg *config.Config) (*AnalyzerService, error) {
	// 创建用户专属的输出目录
	dir_name := cfg.UID + "_" + cfg.OutputName
	userOutputDir := filepath.Join(cfg.OutputDir, dir_name)
	if err := os.MkdirAll(userOutputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建文件夹失败%w", err)
	}

	// 创建统计数据文件
	statsFilePath := filepath.Join(userOutputDir, "stats.txt")
	statsFileLLMPath := filepath.Join(userOutputDir, "stats-llm.cvs")
	statsFile, err := os.OpenFile(statsFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	statsFileLLM, err := os.OpenFile(statsFileLLMPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return nil, fmt.Errorf("创建统计数据文件失败: %v\n", err)
	}

	return &AnalyzerService{
		weiboService: weiboService,
		statistics: &models.PhoneStatistics{
			BrandCounts: make(map[string]int),
			UserCount:   0,
		},
		processedUsers: make(map[string]bool),
		outputDir:      userOutputDir,
		statsFile:      statsFile,
		statsFileLLM:   statsFileLLM,
		interval:       cfg.Interval,
		LLMAnalysis: llm.DeepSeek{
			Api_key: cfg.ApiKey,
		},
	}, nil
}

// AnalyzeUserPhones 分析用户手机品牌分布
func (a *AnalyzerService) AnalyzeUserPhones(uid string, blog_list []string) *models.PhoneStatistics {
	fmt.Printf("开始分析用户 %s 的手机品牌分布，限制 %d 个用户\n", uid)

	// 重置统计
	a.resetStatistics()

	// 定义用户处理回调
	userCallback := func(comments []models.CommentData, blog_content string) {
		a.processUsers(comments, blog_content, a.interval)
	}

	// 获取并处理用户
	a.weiboService.GetUserBlogsAndComments(uid, blog_list, a.interval, userCallback)

	fmt.Printf("分析完成，共处理 %d 个用户\n", a.statistics.UserCount)
	return a.statistics
}

// isUserProcessed 检查用户是否已处理
func (a *AnalyzerService) isUserProcessed(userID string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.processedUsers[userID]
}

// markUserAsProcessed 标记用户为已处理
func (a *AnalyzerService) markUserAsProcessed(userID string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.processedUsers[userID] = true
}

// processUsers 处理用户列表
func (a *AnalyzerService) processUsers(comments []models.CommentData, blog_content string, interval int) {
	for _, comment := range comments {
		// 检查用户是否已处理过（全局去重）
		if a.isUserProcessed(comment.User.ID) {
			continue
		}

		// 获取用户手机类型
		phoneType, err := a.weiboService.GetUserPhoneType(comment.User.ID)
		if err != nil {
			fmt.Printf("获取用户 %s 手机类型失败: %v，跳过\n", comment.User.ID, err)
			continue
		}

		// 标记用户为已处理
		a.markUserAsProcessed(comment.User.ID)
		// 实时写入用户统计数据到文件
		a.writeUserStats(comment.User.ID, phoneType)
		// 更新统计
		a.updateStatistics(phoneType)

		go func() {
			result, err := a.LLMAnalysis.GetCommentLevel(comment.Text, blog_content)
			if err != nil {
				fmt.Println("调用LLM API 发生错误 %v", err)
				return
			}
			a.statisticsLLM[comment.User.ID] = models.StatisticData{
				ResonContent: result.ReasonContent,
				UID:          comment.User.ID,
				Value:        result.Value,
				PhoneType:    phoneType,
			}
		}()
		// 避免请求过于频繁
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// updateStatistics 更新统计信息
func (a *AnalyzerService) updateStatistics(phoneType string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.statistics.BrandCounts[phoneType]++
	a.statistics.UserCount++
}

// writeUserStats 实时写入用户统计数据到文件
func (a *AnalyzerService) writeUserStats(userID, phoneType string) {
	if a.statsFile == nil {
		return
	}

	// 实时写入用户ID和设备信息
	statsLine := fmt.Sprintf("%s:%s\n", userID, phoneType)
	if _, err := a.statsFile.WriteString(statsLine); err != nil {
		fmt.Printf("写入统计数据失败: %v\n", err)
	}

	// 立即刷新到磁盘，确保数据不丢失
	if err := a.statsFile.Sync(); err != nil {
		fmt.Printf("刷新统计数据文件失败: %v\n", err)
	}
}

// resetStatistics 重置统计信息
func (a *AnalyzerService) resetStatistics() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.statistics.BrandCounts = make(map[string]int)
	a.statistics.UserCount = 0
	a.processedUsers = make(map[string]bool) // 重置已处理用户集合

	// 重置统计数据文件
	if a.statsFile != nil {
		a.statsFile.Close()
		statsFilePath := filepath.Join(a.outputDir, "stats.txt")
		statsFile, err := os.OpenFile(statsFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("重置统计数据文件失败: %v\n", err)
			a.statsFile = nil
		} else {
			a.statsFile = statsFile
		}
	}
}

// GetStatistics 获取统计信息
func (a *AnalyzerService) GetStatistics() *models.PhoneStatistics {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// 创建副本以避免并发问题
	statistics := &models.PhoneStatistics{
		BrandCounts: make(map[string]int),
		UserCount:   a.statistics.UserCount,
	}

	for k, v := range a.statistics.BrandCounts {
		statistics.BrandCounts[k] = v
	}

	return statistics
}

// GetKnownBrandStats 获取已知品牌统计
func (a *AnalyzerService) GetKnownBrandStats() []models.StatisticsData {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var result []models.StatisticsData
	for phoneType, count := range a.statistics.BrandCounts {
		if a.weiboService.IsKnownBrand(phoneType) {
			result = append(result, models.StatisticsData{
				PhoneType: phoneType,
				Count:     count,
			})
		}
	}

	// 按数量排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}

// GetUnknownBrandStats 获取未知品牌统计
func (a *AnalyzerService) GetUnknownBrandStats() []models.StatisticsData {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var result []models.StatisticsData
	for phoneType, count := range a.statistics.BrandCounts {
		if !a.weiboService.IsKnownBrand(phoneType) {
			result = append(result, models.StatisticsData{
				PhoneType: phoneType,
				Count:     count,
			})
		}
	}

	// 按数量排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}

// PrintProgress 打印当前进度
func (a *AnalyzerService) PrintProgress() {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	fmt.Printf("当前进度: 已处理 %d 个用户\n", a.statistics.UserCount)

	// 打印已知品牌统计
	knownStats := a.GetKnownBrandStats()
	fmt.Println("已知品牌统计:")
	for _, stat := range knownStats {
		fmt.Printf("  %s: %d\n", stat.PhoneType, stat.Count)
	}
	fmt.Println(strings.Repeat("-", 50))
}

// GetProcessedUserCount 获取已处理用户数量（包括去重统计）
func (a *AnalyzerService) GetProcessedUserCount() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return len(a.processedUsers)
}

// GetDuplicateUserCount 获取重复用户数量（用于统计去重效果）
func (a *AnalyzerService) GetDuplicateUserCount() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	// 重复用户数 = 总处理用户数 - 唯一用户数
	return a.statistics.UserCount - len(a.processedUsers)
}

// GetSummary 获取分析摘要
func (a *AnalyzerService) GetSummary() string {
	stats := a.GetStatistics()
	knownStats := a.GetKnownBrandStats()
	unknownStats := a.GetUnknownBrandStats()
	uniqueUserCount := a.GetProcessedUserCount()
	duplicateCount := a.GetDuplicateUserCount()

	var builder strings.Builder
	builder.WriteString("=== 分析摘要 ===\n")
	builder.WriteString(fmt.Sprintf("唯一用户数: %d\n", uniqueUserCount))
	if duplicateCount > 0 {
		builder.WriteString(fmt.Sprintf("重复用户数: %d\n", duplicateCount))
	}
	builder.WriteString(fmt.Sprintf("总处理用户数: %d\n", stats.UserCount))
	builder.WriteString(fmt.Sprintf("已知品牌数: %d\n", len(knownStats)))
	builder.WriteString(fmt.Sprintf("未知品牌数: %d\n", len(unknownStats)))

	if len(knownStats) > 0 {
		builder.WriteString("\n前5名已知品牌:\n")
		for i, stat := range knownStats {
			if i >= 5 {
				break
			}
			builder.WriteString(fmt.Sprintf("  %d. %s: %d (%.1f%%)\n",
				i+1, stat.PhoneType, stat.Count,
				float64(stat.Count)/float64(uniqueUserCount)*100))
		}
	}

	if len(unknownStats) > 0 {
		builder.WriteString(fmt.Sprintf("\n未知品牌数量: %d\n", len(unknownStats)))
		if len(unknownStats) <= 10 {
			builder.WriteString("未知品牌列表:\n")
			for _, stat := range unknownStats {
				builder.WriteString(fmt.Sprintf("  %s: %d\n", stat.PhoneType, stat.Count))
			}
		}
	}

	return builder.String()
}

// Close 关闭分析服务，释放资源
func (a *AnalyzerService) Close() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.statsFile != nil {
		err := a.statsFile.Close()
		a.statsFile = nil
		return err
	}
	return nil
}

// GetOutputDir 获取输出目录路径
func (a *AnalyzerService) GetOutputDir() string {
	return a.outputDir
}
