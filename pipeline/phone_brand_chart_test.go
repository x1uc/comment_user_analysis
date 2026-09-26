package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/x1uc/comment_user_analysis/models"
)

func TestPhoneBrandChartCountsAndWritesHTML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "phone_brand_pie.html")
	plugin := NewPhoneBrandChartPlugin(path)
	mem := &Memory{
		PhoneInfos: []models.UserPhoneInfo{
			{PhoneBrand: "华为"},
			{PhoneBrand: "华为"},
			{PhoneBrand: "小米"},
			{PhoneBrand: "  "},
			{PhoneBrand: "Android设备"},
		},
	}
	if err := plugin.Run(context.Background(), mem); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	html := string(body)
	for _, want := range []string{
		`["华为","未知设备","小米"]`,
		`[2,2,1]`,
		`["#CF0A2C","#A9A9A9","#FF6900"]`,
		"chartXkcd.Pie",
		"showLegend: true",
		"drawBrandCallouts",
		"brand-arrow",
		"https://cdn.jsdelivr.net/npm/chart.xkcd@1/dist/chart.xkcd.min.js",
		"innerRadius: 0",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("html missing %s\n%s", want, html)
		}
	}
}

func TestPhoneBrandChartEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "output", "phone_brand_pie.html")
	if err := NewPhoneBrandChartPlugin(path).Run(context.Background(), &Memory{}); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	html := string(body)
	if !strings.Contains(html, "labels: []") || !strings.Contains(html, "data: []") {
		t.Fatalf("empty chart html = %s", html)
	}
}
