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
		Model:     anthropic.ModelClaude3_5Sonnet20241022,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}

	if len(message.Content) > 0 {
		if textBlock, ok := message.Content[0].(anthropic.TextBlock); ok {
			return textBlock.Text, nil
		}
	}

	return "", fmt.Errorf("no text content in response")
}
