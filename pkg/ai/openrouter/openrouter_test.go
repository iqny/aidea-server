package openrouter_test

import (
	"context"
	"encoding/json"
	openai2 "github.com/mylxsw/aidea-server/pkg/ai/openai"
	"github.com/mylxsw/aidea-server/pkg/ai/openrouter"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-utils/assert"
	"github.com/mylxsw/go-utils/must"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestOpenRouter_ChatStream(t *testing.T) {
	client := createClient()

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	response, err := client.ChatStream(ctx, openai.ChatCompletionRequest{
		Model: "deepseek/deepseek-r1",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "9.9和9.11谁大？",
			},
		},
		MaxTokens: 500,
		Stop:      []string{"</think>"},
	})

	assert.NoError(t, err)

	var finalText string
	var reasoning string
	for res := range response {
		print(string(must.Must(json.Marshal(res))))
		if res.Code != "" {
			log.With(res).Errorf("-> %s", res.ErrorMessage)
			continue
		}

		log.Debugf("-> %s", res.ChatResponse.Choices[0].Delta.Content)
		finalText += res.ChatResponse.Choices[0].Delta.Content
		reasoning += res.ChatResponse.Choices[0].Delta.Reasoning
	}

	log.Debugf("reasoning: %s", reasoning)
	log.Debugf("final text: %s", finalText)
}

func createClient() *openrouter.OpenRouter {
	conf := &openai2.Config{
		Enable:        true,
		OpenAIServers: []string{"https://openrouter.ai/api/v1"},
		OpenAIKeys:    []string{os.Getenv("OPEN_ROUTER_API_KEY")},
		Header: http.Header{
			"HTTP-Referer": []string{"https://web.aicode.cc"},
			"X-Title":      []string{"AIdea"},
		},
	}
	return openrouter.NewOpenRouter(openai2.NewOpenAIClient(conf, nil))
}

func TestOpenRouter_Chat(t *testing.T) {
	client := createClient()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	response, err := client.Chat(ctx, openai.ChatCompletionRequest{
		Model: "openai/o3-mini-high",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "假如你是鲁迅，请使用批判性，略带讽刺的语言来回答我的问题，语言要风趣，幽默，略带调侃",
			},
			{
				Role:    "user",
				Content: "老铁，最近怎么样？",
			},
		},
	})

	assert.NoError(t, err)

	log.With(response).Debugf("response")
}
