package llm

import (
	"context"
	"fmt"
	"single_analysis/internal/models"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var (
	test_prompt string = `
						# Role你是一位精通中国互联网舆论生态、法学辩论逻辑以及中文网络俚语（包括反讽、隐喻）的社会舆论分析师。

						# Background
						背景：近期中国互联网上关于“吸毒人员违法犯罪记录封存”制度引发巨大争议。
						核心人物：劳东燕（清华大学法学院教授）。
						人物立场：劳东燕支持建立轻罪/特定条件下吸毒记录封存制度，主张给改过自新者回归社会的机会。
						舆论环境：大众普遍对毒品持“零容忍”态度，常引用缉毒警的牺牲来反对任何形式的宽容。评论区可能包含激烈的对立情绪。

						# Task
						请分析我提供的微博评论列表，将每一条评论分类为以下三类之一：
						1. 【支持封存】Value 2：赞同劳东燕的观点。关键词包括：人权、回归社会、法治精神、给条活路、反对连坐等。
						2. 【反对封存】Value 1：反对劳东燕的观点，反对消除吸毒记录。关键词包括：零容忍、缉毒警、牺牲、圣母、复吸、对受害者不公、何不食肉糜等。
						*注意：如果评论是在人身攻击劳东燕（如骂她“公知”、“西方法律代理人”），在语境下也应视为“反对封存”的立场。
						3. 【中立/其他】Value 0：纯吃瓜、或讨论完全无关的话题、评论信息不足以分辨立场（如提及安卓苹果大战但未与本题建立逻辑联系）。

						# Constraints & Rules
						- **识别反讽**：这是最关键的。如果用户说“教授真是大善人，建议把吸毒者领回自己家养”，表面是夸赞，实际是极度【反对封存】。请务必识别“友军厚葬”类的反串言论。
						- **逻辑关联**：如果评论提及“安卓/苹果”话题，除非用户用其进行类比（如“保护毒虫隐私却不保护苹果用户隐私”），否则归为【中立/其他】。
						- **输出格式**：请以JSON格式输出，例如 {Value : 1}。

						# 劳动燕微博正文内容
						%s
						`
)

type LLMClient interface {
	GetCommentLevel(comment string, blog_content string) (*models.AiResponse, error)
}

type DeepSeek struct {
	Api_key string
}

func (c DeepSeek) GetCommentLevel(comment string, blog_content string) (*models.AiResponse, error) {
	client := openai.NewClient(
		option.WithBaseURL("https://api.deepseek.com"),
		option.WithAPIKey(c.Api_key), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel("deepseek-reasoner"),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(fmt.Sprintf(test_prompt, blog_content)),
			openai.UserMessage(comment),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{
				Type: "json_object",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	choice := chatCompletion.Choices[0]

	var reason_content string
	if reason_content, ok := choice.Message.JSON.ExtraFields["reasoning_content"]; ok {
		fmt.Printf("=== 思考过程 ===\n%v\n\n", reason_content)
	}
	res := &models.AiResponse{
		Value:         100,
		ReasonContent: reason_content,
	}
	return res, nil
}
