# 微博用户手机品牌分析工具

一个用于分析微博用户手机品牌分布的工具，支持统计用户评论者的手机品牌并生成可视化图表。

## 使用方法

创建 `config.json` 文件：

```json
{
  "uid": "需要的用户UID",
  "cookie": "你的微博Cookie",
  "limit": 100,
  "debug": false,
  "output_dir": "./output"
}
```

运行 
```
go mod tidy

go run cmd/main.go
```

## FAQ

如果不能使用，请修改 intermal/client/client.go 中的 setHeaders ，保证和当前微博网页端同步