package controllers

import (
    "context"
    "log"
	"fmt"
	"time"
    "github.com/gin-gonic/gin"
    openai "github.com/sashabaranov/go-openai"
    "os"
)

type ChatRequest struct {
    Message string `json:"message" binding:"required"`
}

func HandleChat(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    response, err := ChatGenerate(req.Message)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"response": response})
}

func ChatGenerate(userInput string) (string, error) {
	apiKey := os.Getenv("ChatGPT_SECRET_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ChatGPT_SECRET_KEY not set")
	}
	
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	Today := time.Now().Format("2006년 01월 02일 Monday")
	client := openai.NewClient(apiKey)
	baseStr := "현재 날짜와 시간이" + Today + currentTime + "에 9시간을 더한 것인데 날짜로 한국 계절을 반영하고 현재 시간을 반영해서 메뉴 추천을 해줄 수 있을까? 답변에 시간은 포함하지 말고."
	
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userInput + baseStr,
				},
			},
		},
	)
	
	if err != nil {
		log.Printf("ChatGPT API Error: %v", err)
		return "", err
	}
	
	return resp.Choices[0].Message.Content, nil
 }