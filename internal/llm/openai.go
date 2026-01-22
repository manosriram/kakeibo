package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client *openai.Client
}

func NewOpenAI() OpenAI {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	return OpenAI{
		Client: client,
	}
}

func (o OpenAI) GeneratePrompt() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	d, err := os.ReadFile(wd + "/internal/llm/prompt.md")
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func (o OpenAI) Call(p string) (string, error) {
	prompt, err := o.GeneratePrompt()
	if err != nil {
		return "", err
	}
	prompt = fmt.Sprintf(prompt, p)
	chatCompletion, err := o.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT5,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func (o OpenAI) CreateEmbedding(payload map[string]any) ([]openai.Embedding, error) {
	var ip strings.Builder
	for k, v := range payload {
		ip.WriteString(fmt.Sprintf("%s:%s", k, v))
	}
	req := openai.EmbeddingRequest{
		Model: openai.SmallEmbedding3,
		Input: ip.String(),
	}
	embeddingResponse, err := o.Client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		return []openai.Embedding{}, err
	}

	return embeddingResponse.Data, nil

}
