package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
)

const (
	quitMessage = "quit"
)

func getResponse(client *openai.Client, ctx context.Context, message string) (*string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		}},
		Stream: false,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "ChatCompletion error")
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("ChatCompletion request returned no responses")
	}

	return &resp.Choices[0].Message.Content, nil
}

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("missing api key")
		os.Exit(1)
	}

	client := openai.NewClient(apiKey)

	scanner := bufio.NewScanner(os.Stdin)

	quit := false

	for !quit {
		fmt.Print("Input your message (type `quit` to exit): ")

		if !scanner.Scan() {
			break
		}

		msg := strings.Trim(scanner.Text(), " ")
		switch msg {
		case quitMessage:
			quit = true
		default:
			resp, err := getResponse(client, ctx, msg)
			if err != nil {
				fmt.Printf("error: %s", err)
				os.Exit(1)
			}
			fmt.Println(*resp)
		}
	}
}
