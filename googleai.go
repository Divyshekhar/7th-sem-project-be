package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

func callGemini() {
	ctx := context.Background()

	// apiKey := os.Getenv("GEMINI_API")
	llm, err := googleai.New(ctx, googleai.WithAPIKey(""), googleai.WithDefaultModel("gemini-2.0-flash"))
	if err != nil {
		log.Fatal(err)
	}

	// for streaming messages

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(llms.ChatMessageTypeHuman, "What would be a good company name for a comapny that produces Go-backed LLM tools?"),
	}

	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
}
