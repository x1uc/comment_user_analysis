package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"single_analysis/config"
	"single_analysis/export"
	"single_analysis/internal/models"
	"single_analysis/internal/services"
	"syscall"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 打印配置信息
	cfg.Print()

	// 创建服务
	weiboService := services.NewWeiboService(cfg.Cookie)
	analyzerService, err := services.NewAnalyzerService(weiboService, cfg)
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer analyzerService.Close() // 确保资源释放

	setupGracefulShutdown(analyzerService, cfg)

	// 开始分析
	fmt.Println("开始分析...")
	analyzerService.AnalyzeUserPhones(cfg.UID, cfg.BlogList, cfg.Interval)

	// 打印结果
	printResults(analyzerService)
	convertDataToChart(analyzerService, cfg)
}

func convertDataToChart(analyzerService *services.AnalyzerService, cfg *config.Config) {
	// 导出图表到用户专属目录
	userOutputDir := analyzerService.GetOutputDir()
	chartExporter := export.NewChartExporter(cfg.UID, userOutputDir)
	fmt.Println("\n开始导出图表...")
	knownStats := analyzerService.GetKnownBrandStats()

	// 导出饼图
	if err := chartExporter.ExportPieChart(knownStats); err != nil {
		log.Printf("导出饼图失败: %v", err)
	} else {
		fmt.Println("饼图导出完成!")
	}

	// 导出柱状图
	if err := chartExporter.ExportBarChart(knownStats); err != nil {
		log.Printf("导出柱状图失败: %v", err)
	} else {
		fmt.Println("柱状图导出完成!")
	}

	// 导出摘要
	summaryExporter := export.NewChartExporter(cfg.UID, userOutputDir)

	// 获取所有统计数据（包括未知机型）
	allStats := analyzerService.GetStatistics()
	var allStatsData []models.StatisticsData
	for phoneType, count := range allStats.BrandCounts {
		allStatsData = append(allStatsData, models.StatisticsData{
			PhoneType: phoneType,
			Count:     count,
		})
	}

	if err := summaryExporter.ExportSummary(allStatsData); err != nil {
		log.Printf("导出摘要失败: %v", err)
	} else {
		fmt.Println("摘要导出完成!")
	}

	// 打印摘要
	fmt.Println("\n" + analyzerService.GetSummary())
	fmt.Printf("所有文件已保存到目录: %s\n", userOutputDir)
}

// printResults 打印分析结果
func printResults(analyzer *services.AnalyzerService) {
	fmt.Println("\n========================== 最终统计结果：未知机型 ============================")
	unknownStats := analyzer.GetUnknownBrandStats()
	for _, stat := range unknownStats {
		fmt.Printf("PhoneType: %s, Num: %d\n", stat.PhoneType, stat.Count)
	}

	fmt.Println("\n========================== 最终统计结果：已知机型 ============================")
	knownStats := analyzer.GetKnownBrandStats()
	for _, stat := range knownStats {
		fmt.Printf("PhoneType: %s, Num: %d\n", stat.PhoneType, stat.Count)
	}

	// 计算总数
	totalKnown := 0
	for _, stat := range knownStats {
		totalKnown += stat.Count
	}
	fmt.Printf("\n已知品牌总用户数: %d\n", totalKnown)
}

// setupGracefulShutdown 设置优雅退出
func setupGracefulShutdown(analyzerService *services.AnalyzerService, cfg *config.Config) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n\n收到退出信号，正在优雅退出...")
		convertDataToChart(analyzerService, cfg)
		os.Exit(0)
	}()
}
