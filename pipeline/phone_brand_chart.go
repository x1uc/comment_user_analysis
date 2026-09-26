package pipeline

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/x1uc/comment_user_analysis/models"
)

const unknownPhoneBrand = "未知设备"

type PhoneBrandChartPlugin struct {
	path string
}

func NewPhoneBrandChartPlugin(path string) *PhoneBrandChartPlugin {
	return &PhoneBrandChartPlugin{path: path}
}

func (p *PhoneBrandChartPlugin) Name() string {
	return "phone_brand_chart"
}

func (p *PhoneBrandChartPlugin) Run(ctx context.Context, mem *Memory) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	labels, counts := countPhoneBrands(mem.PhoneInfos)
	if err := writePhoneBrandChart(p.path, labels, counts, brandChartColors(labels)); err != nil {
		return err
	}
	log.Printf("phone brand chart: %s", p.path)
	return nil
}

func countPhoneBrands(infos []models.UserPhoneInfo) ([]string, []int) {
	countsByBrand := make(map[string]int)
	for _, info := range infos {
		countsByBrand[chartPhoneBrand(info.PhoneBrand)]++
	}
	labels := make([]string, 0, len(countsByBrand))
	for brand := range countsByBrand {
		labels = append(labels, brand)
	}
	sort.Slice(labels, func(i, j int) bool {
		if countsByBrand[labels[i]] == countsByBrand[labels[j]] {
			return labels[i] < labels[j]
		}
		return countsByBrand[labels[i]] > countsByBrand[labels[j]]
	})
	counts := make([]int, len(labels))
	for i, label := range labels {
		counts[i] = countsByBrand[label]
	}
	return labels, counts
}

func brandChartColors(labels []string) []string {
	colors := make([]string, len(labels))
	for i, label := range labels {
		colors[i] = models.BrandColor(label)
	}
	return colors
}

func chartPhoneBrand(brand string) string {
	switch strings.ToLower(strings.TrimSpace(brand)) {
	case "", "未知", "android设备", "安卓设备", "android":
		return unknownPhoneBrand
	default:
		return strings.TrimSpace(brand)
	}
}

func writePhoneBrandChart(path string, labels []string, counts []int, colors []string) error {
	if labels == nil {
		labels = []string{}
	}
	if counts == nil {
		counts = []int{}
	}
	if colors == nil {
		colors = []string{}
	}
	labelsJSON, err := json.Marshal(labels)
	if err != nil {
		return err
	}
	countsJSON, err := json.Marshal(counts)
	if err != nil {
		return err
	}
	colorsJSON, err := json.Marshal(colors)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return phoneBrandChartTemplate.Execute(file, map[string]template.JS{
		"Labels": template.JS(labelsJSON),
		"Data":   template.JS(countsJSON),
		"Colors": template.JS(colorsJSON),
	})
}

var phoneBrandChartTemplate = template.Must(template.New("phone_brand_chart").Parse(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <title>手机品牌</title>
</head>
<body>
  <svg class="pie-chart"></svg>
  <script src="https://cdn.jsdelivr.net/npm/chart.xkcd@1/dist/chart.xkcd.min.js"></script>
  <script>
    const svg = document.querySelector('.pie-chart');
    new chartXkcd.Pie(svg, {
      title: '手机品牌',
      data: {
        labels: {{.Labels}},
        datasets: [{
          data: {{.Data}},
        }],
      },
      options: {
        innerRadius: 0,
        showLegend: true,
        dataColors: {{.Colors}},
      },
    });
    drawBrandCallouts(svg, {{.Labels}}, {{.Data}}, {{.Colors}});

    function drawBrandCallouts(svg, labels, values, colors) {
      const total = values.reduce(function(sum, value) { return sum + value; }, 0);
      if (!total) {
        return;
      }
      const width = svg.width.baseVal.value;
      const height = svg.height.baseVal.value;
      const pad = 150;
      svg.setAttribute('viewBox', [-pad, -pad, width + pad * 2, height + pad * 2].join(' '));
      const cx = width / 2;
      const cy = height / 2;
      const radius = Math.min(width, height) / 2 - 50;
      const slices = labels.map(function(label, index) {
        return {label: label, value: values[index], index: index};
      });
      slices.sort(function(a, b) {
        return b.value - a.value || a.index - b.index;
      });
      const left = [];
      const right = [];
      let angle = -Math.PI / 2;
      slices.forEach(function(slice) {
        const sweep = slice.value / total * Math.PI * 2;
        const mid = angle + sweep / 2;
        angle += sweep;
        const cos = Math.cos(mid);
        const sin = Math.sin(mid);
        const item = {
          label: slice.label,
          x1: cx + cos * radius,
          y1: cy + sin * radius,
          x2: cx + cos * (radius + 16),
          y2: cy + sin * (radius + 16),
        };
        item.y = item.y2;
        (cos >= 0 ? right : left).push(item);
      });
      const ns = 'http://www.w3.org/2000/svg';
      const defs = document.createElementNS(ns, 'defs');
      defs.innerHTML = '<marker id="brand-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill="#222"></path></marker>';
      svg.insertBefore(defs, svg.firstChild);
      placeCallouts(svg, ns, right, cx + radius + 28, 1, -pad + 18, height + pad - 18);
      placeCallouts(svg, ns, left, cx - radius - 28, -1, -pad + 18, height + pad - 18);
    }

    function placeCallouts(svg, ns, items, labelX, direction, minY, maxY) {
      items.sort(function(a, b) { return a.y - b.y; });
      for (let i = 1; i < items.length; i++) {
        if (items[i].y < items[i - 1].y + 22) {
          items[i].y = items[i - 1].y + 22;
        }
      }
      if (items.length && items[items.length - 1].y > maxY) {
        const shift = items[items.length - 1].y - maxY;
        items.forEach(function(item) { item.y -= shift; });
      }
      if (items.length && items[0].y < minY) {
        const shift = minY - items[0].y;
        items.forEach(function(item) { item.y += shift; });
      }
      items.forEach(function(item) {
        const line = document.createElementNS(ns, 'polyline');
        line.setAttribute('points', [item.x1, item.y1, item.x2, item.y2, labelX, item.y].join(' '));
        line.setAttribute('fill', 'none');
        line.setAttribute('stroke', '#222');
        line.setAttribute('stroke-width', '1.5');
        line.setAttribute('marker-end', 'url(#brand-arrow)');
        svg.appendChild(line);
        const text = document.createElementNS(ns, 'text');
        text.setAttribute('x', labelX + direction * 8);
        text.setAttribute('y', item.y);
        text.setAttribute('dominant-baseline', 'middle');
        text.setAttribute('text-anchor', direction > 0 ? 'start' : 'end');
        text.setAttribute('font-size', '16');
        text.textContent = item.label;
        svg.appendChild(text);
      });
    }
  </script>
</body>
</html>
`))
