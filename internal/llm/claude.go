package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Claude struct {
	ExpenseDescription string
}

func NewClaude(description string) Claude {
	return Claude{
		ExpenseDescription: description,
	}
}

func (c Claude) GeneratePrompt() (string, error) {
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

func (c Claude) Call() (string, error) {
	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)

	prompt, err := c.GeneratePrompt()
	if err != nil {
		return "", err
	}
	prompt = fmt.Sprintf(prompt, c.ExpenseDescription)

	message, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			{
				Role: anthropic.F(anthropic.MessageParamRoleUser),
				Content: anthropic.F([]anthropic.ContentBlockParamUnion{
					anthropic.TextBlockParam{
						Type: anthropic.F(anthropic.TextBlockParamTypeText),
						Text: anthropic.F(prompt),
					},
				}),
			},
		}),
	})
	if err != nil {
		return "", err
	}

	if len(message.Content) > 0 {
		block := message.Content[0]
		if block.Type == anthropic.ContentBlockTypeText {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("no text content in response")
}
