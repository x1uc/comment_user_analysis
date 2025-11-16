package export

import (
	"comment_phone_analyse/internal/models"
	"comment_phone_analyse/internal/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// ChartExporter 图表导出器
type ChartExporter struct {
	outputDir string
	uid       string
}

// NewChartExporter 创建图表导出器
func NewChartExporter(uid, outputDir string) *ChartExporter {
	return &ChartExporter{
		outputDir: outputDir,
		uid:       uid,
	}
}

// ExportAll 导出所有图表
func (e *ChartExporter) ExportAll(data []models.StatisticsData) error {
	if err := e.ExportBarChart(data); err != nil {
		return utils.NewExportError("导出柱状图失败", err)
	}

	if err := e.ExportPieChart(data); err != nil {
		return utils.NewExportError("导出饼图失败", err)
	}

	if err := e.ExportSummary(data); err != nil {
		return utils.NewExportError("导出摘要失败", err)
	}

	return nil
}

// ExportBarChart 导出柱状图
func (e *ChartExporter) ExportBarChart(data []models.StatisticsData) error {
	if len(data) == 0 {
		return utils.NewExportError("没有数据可导出", nil)
	}

	bar := charts.NewBar()

	// 准备数据
	var xLabels []string
	var yValues []opts.BarData

	for _, phone := range data {
		xLabels = append(xLabels, phone.PhoneType)
		yValues = append(yValues, opts.BarData{
			Value:     phone.Count,
			ItemStyle: &opts.ItemStyle{Color: e.getColor(phone.PhoneType)},
			Label: &opts.Label{
				Show:      &[]bool{true}[0],
				Position:  "top",
				Formatter: "{c}",
			},
		})
	}

	// 设置全局选项
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("用户 %s 的手机品牌分布", e.uid),
			Subtitle: "柱状图统计",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "手机品牌",
			AxisLabel: &opts.AxisLabel{
				Interval: strconv.Itoa(0),
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{Name: "用户数量"}),
		charts.WithGridOpts(opts.Grid{
			Left:   "10%",
			Right:  "10%",
			Bottom: "15%",
			Top:    "15%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: fmt.Sprintf("%s 手机品牌统计", e.uid),
		}),
	)

	bar.SetXAxis(xLabels).AddSeries("用户数量", yValues)

	// 保存文件
	filename := filepath.Join(e.outputDir, "stats.html")
	return e.saveChart(bar, filename)
}

// ExportPieChart 导出饼图
func (e *ChartExporter) ExportPieChart(data []models.StatisticsData) error {
	if len(data) == 0 {
		return utils.NewExportError("没有数据可导出", nil)
	}

	pie := charts.NewPie()

	// 准备数据
	var pieData []opts.PieData
	for _, phone := range data {
		pieData = append(pieData, opts.PieData{
			Name:      phone.PhoneType,
			Value:     phone.Count,
			ItemStyle: &opts.ItemStyle{Color: e.getColor(phone.PhoneType)},
			Label: &opts.Label{
				Show:      &[]bool{true}[0],
				Formatter: "{b}: {d}%",
			},
		})
	}

	// 设置全局选项
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("用户 %s 的手机品牌分布", e.uid),
			Subtitle: "饼图统计",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: fmt.Sprintf("%s 手机品牌分布", e.uid),
		}),
	)

	pie.AddSeries("手机品牌", pieData)

	// 保存文件
	filename := filepath.Join(e.outputDir, "pie.html")
	return e.saveChart(pie, filename)
}

// ExportSummary 导出统计摘要
func (e *ChartExporter) ExportSummary(data []models.StatisticsData) error {
	if len(data) == 0 {
		return utils.NewExportError("没有数据可导出", nil)
	}

	filename := filepath.Join(e.outputDir, "summary.txt")
	file, err := os.Create(filename)
	if err != nil {
		return utils.NewExportError("创建摘要文件失败", err)
	}
	defer file.Close()

	// 写入摘要内容
	fmt.Fprintf(file, "=== 用户 %s 手机品牌统计摘要 ===\n\n", e.uid)

	total := 0
	for _, phone := range data {
		total += phone.Count
	}

	fmt.Fprintf(file, "统计时间: %s\n", getNowTime())
	fmt.Fprintf(file, "总用户数: %d\n", total)
	fmt.Fprintf(file, "品牌数量: %d\n\n", len(data))

	fmt.Fprintf(file, "详细统计:\n")
	// 按数量降序排序
	sort.Slice(data, func(i, j int) bool {
		return data[i].Count > data[j].Count
	})
	for i, phone := range data {
		percentage := float64(phone.Count) / float64(total) * 100
		fmt.Fprintf(file, "%2d. %-12s: %4d (%5.1f%%)\n", i+1, phone.PhoneType, phone.Count, percentage)
	}

	// 添加未知机型信息
	fmt.Fprintf(file, "\n========================== 未知机型统计 ===========================\n")
	unknownCount := 0
	knownBrands := map[string]bool{
		"华为": true, "小米": true, "OPPO": true, "Vivo": true, "苹果": true,
		"三星": true, "魅族": true, "真我": true, "红米": true, "一加": true,
		"荣耀": true, "中兴": true, "努比亚": true, "IQOO": true, "未知Android": true,
	}

	for _, phone := range data {
		if !knownBrands[phone.PhoneType] {
			fmt.Fprintf(file, "PhoneType: %s, Num: %d\n", phone.PhoneType, phone.Count)
			unknownCount += phone.Count
		}
	}

	fmt.Printf("统计摘要已保存到: %s\n", filename)
	return nil
}

// saveChart 保存图表到文件
func (e *ChartExporter) saveChart(chart interface{ Render(io.Writer) error }, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return utils.NewExportError("创建图表文件失败", err)
	}
	defer file.Close()

	if err := chart.Render(file); err != nil {
		return utils.NewExportError("渲染图表失败", err)
	}

	fmt.Printf("图表已保存到: %s\n", filename)
	return nil
}

// getColor 获取品牌对应的颜色
func (e *ChartExporter) getColor(phoneType string) string {
	colorsMap := map[string]string{
		"真我":        "#FFD700",
		"苹果":        "#333333",
		"OPPO":      "#008000",
		"荣耀":        "#0033A0",
		"华为":        "#FF0000",
		"Vivo":      "#0072C6",
		"一加":        "#FF0000",
		"小米":        "#FF6700",
		"红米":        "#ED1C24",
		"IQOO":      "#FF6F00",
		"魅族":        "#C71585",
		"三星":        "#1428A0",
		"努比亚":       "#FF0000",
		"未知Android": "#A9A9A9",
		"Other":     "#A9954B",
	}

	if color, exists := colorsMap[phoneType]; exists {
		return color
	}
	return colorsMap["Other"]
}

// getNowTime 获取当前时间字符串
func getNowTime() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second())
}
