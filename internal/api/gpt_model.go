package api

import (
	"context"
	"fmt"
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func SendPromptToModel(userPrompt string, systemPrompt string) (string, error) {
	client := openai.NewClient(option.WithAPIKey(config.OpenAICfg.SecretKey))

	ctx := context.Background()

	// 프롬프트 작성
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "gpt-4.1-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
	})
	if err != nil {
		log.Fatalf("API request error: %v", err)
		return "", err
	}

	fmt.Println(resp.Choices[0].Message.Content) // check
	return resp.Choices[0].Message.Content, nil
}
