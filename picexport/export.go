package picexport

import (
	"comment_phone_analyse/pojo"
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Export(uid string, phoneArr []pojo.StatisticsData) {
	// 创建一个新的柱状图
	bar := charts.NewBar()

	// 准备数据
	var xLabels []string
	var yValues []opts.BarData
	for _, phone := range phoneArr {
		xLabels = append(xLabels, phone.PhoneType)
		yValues = append(yValues, opts.BarData{
			Value:     phone.Count,
			ItemStyle: &opts.ItemStyle{Color: getColor(phone.PhoneType)},
			Label: &opts.Label{
				Show:      &[]bool{true}[0],
				Position:  "top",
				Formatter: "{c}",
			},
		})
	}

	// 设置柱状图数据
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Phone Type Statistics"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Phone Type"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Number"}),
		charts.WithGridOpts(opts.Grid{Left: "10%", Right: "10%", Bottom: "10%", Top: "10%"}), // 设置图表区域的边距
	)
	bar.SetXAxis(xLabels).AddSeries("Phone Count", yValues)

	// 保存图表为HTML文件
	f, err := os.Create(fmt.Sprintf("%s_phone_stats.html", uid))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	if err := bar.Render(f); err != nil {
		fmt.Println("Error rendering chart:", err)
	}
}

func ExportPieChart(uid string, phoneArr []pojo.StatisticsData) {
	// 创建一个新的饼状图
	pie := charts.NewPie()

	// 准备数据
	var pieData []opts.PieData
	for _, phone := range phoneArr {
		pieData = append(pieData, opts.PieData{
			Name:      phone.PhoneType,
			Value:     phone.Count,
			ItemStyle: &opts.ItemStyle{Color: getColor(phone.PhoneType)},
			Label: &opts.Label{
				Show:      &[]bool{true}[0],
				Formatter: "{b}: {d}%", // 显示百分比
			},
		})
	}

	// 设置饼状图数据
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: ""}),
	)
	pie.AddSeries("Phone Count", pieData)

	// 保存图表为HTML文件
	f, err := os.Create(fmt.Sprintf("%s_phone_pie.html", uid))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	if err := pie.Render(f); err != nil {
		fmt.Println("Error rendering chart:", err)
	}
}

func getColor(phoneType string) string {
	var colorsMap = map[string]string{
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
