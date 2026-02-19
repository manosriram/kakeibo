package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAI struct {
	ExpenseDescription string
}

func NewOpenAI(description string) OpenAI {
	return OpenAI{
		ExpenseDescription: description,
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

func (o OpenAI) Call() (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	prompt, err := o.GeneratePrompt()
	if err != nil {
		return "", err
	}
	prompt = fmt.Sprintf(prompt, o.ExpenseDescription)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModel("openai/gpt-4.2"),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
