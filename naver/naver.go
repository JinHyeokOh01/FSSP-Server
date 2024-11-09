package naver

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
	"os"
    "log"
	"github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
)

type NaverSearchResponse struct {
    Items []struct {
        Title       string `json:"title"`
        Category        string `json:"category"`
        Description string `json:"description"`
        Telephone     string `json:"telephone"`
        RoadAddress string `json:"roadAddress"`
        Mapx        string `json:"mapx"`
        Mapy        string `json:"mapy"`
    } `json:"items"`
}

func NaverSearch(query string, display int) (*NaverSearchResponse, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    ClientID := os.Getenv("Naver_Client_ID")
    ClientSecret := os.Getenv("Naver_Secret")

    // URL 인코딩
    encodedQuery := url.QueryEscape(query)
    url := fmt.Sprintf("https://openapi.naver.com/v1/search/local.json?query=%s&display=%d", encodedQuery, display)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("X-Naver-Client-Id", ClientID)
    req.Header.Add("X-Naver-Client-Secret", ClientSecret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // 상태 코드 확인
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
    }

    // 응답 바디 읽기
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // JSON 디코딩
    var searchResponse NaverSearchResponse
    err = json.Unmarshal(body, &searchResponse)
    if err != nil {
        return nil, err
    }
    
    return &searchResponse, nil
}

func QuerySearch(c *gin.Context) {
    query := c.Query("query") // URL 파라미터에서 쿼리 가져오기
    display := 5
    result, err := NaverSearch(query, display)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}
