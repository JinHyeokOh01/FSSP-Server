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

type ChatResponse struct {
    Response string `json:"response"`
    Error    error
}

func HandleChat(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    // 응답 채널 생성
    responseChan := make(chan ChatResponse, 1)
    
    // 비동기로 ChatGPT API 호출
    go func() {
        response, err := ChatGenerate(req.Message)
        responseChan <- ChatResponse{
            Response: response,
            Error:    err,
        }
    }()

    // 타임아웃과 함께 응답 대기
    select {
    case result := <-responseChan:
        if result.Error != nil {
            c.JSON(500, gin.H{"error": result.Error.Error()})
            return
        }
        c.JSON(200, gin.H{"response": result.Response})
    case <-time.After(30 * time.Second): // ChatGPT API는 응답이 좀 걸릴 수 있으므로 30초 타임아웃
        c.JSON(504, gin.H{"error": "요청 시간이 초과되었습니다"})
        return
    }
}

func ChatGenerate(userInput string) (string, error) {
    apiKey := os.Getenv("ChatGPT_SECRET_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("ChatGPT_SECRET_KEY not set")
    }

    // API 호출 결과를 위한 채널
    resultChan := make(chan string, 1)
    errChan := make(chan error, 1)
    
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    Today := time.Now().Format("2006년 01월 02일 Monday")
    client := openai.NewClient(apiKey)
    baseStr := "현재 날짜와 시간이" + Today + currentTime + "에 9시간을 더한 것인데 날짜로 한국 계절을 반영하고 현재 시간을 반영해서 메뉴 추천을 한 3개 정도 해줄 수 있을까? 식당에 가서 먹을거야. 답변 좀 간결하고 다정하게 해줘. 지금 날짜랑 시간에 대한 언급은 하지마."

    // 컨텍스트 타임아웃 설정
    ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
    defer cancel()

    go func() {
        resp, err := client.CreateChatCompletion(
            ctx,
            openai.ChatCompletionRequest{
                Model: openai.GPT4,
                Messages: []openai.ChatCompletionMessage{
                    {
                        Role:    openai.ChatMessageRoleUser,
                        Content: userInput + baseStr,
                    },
                },
            },
        )

        if err != nil {
            if ctx.Err() == context.DeadlineExceeded {
                errChan <- fmt.Errorf("API 요청 시간이 초과되었습니다")
                return
            }
            log.Printf("ChatGPT API Error: %v", err)
            errChan <- err
            return
        }

        resultChan <- resp.Choices[0].Message.Content
    }()

    // 결과 대기
    select {
    case result := <-resultChan:
        return result, nil
    case err := <-errChan:
        return "", err
    case <-ctx.Done():
        return "", fmt.Errorf("요청이 취소되었습니다: %v", ctx.Err())
    }
}