// searchController.go
package controllers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/JinHyeokOh01/FSSP-Server/db"
)

type NaverSearchResponse struct {
    Items []struct {
        Title       string `json:"title"`
        Category    string `json:"category"`
        Description string `json:"description"`
        Telephone   string `json:"telephone"`
        RoadAddress string `json:"roadAddress"`
        Mapx        string `json:"mapx"`
        Mapy        string `json:"mapy"`
    } `json:"items"`
}

type RestaurantResponse struct {
    Name       string `json:"name"`
    Address    string `json:"address"`
    Category   string `json:"category"`
    IsFavorite bool   `json:"isFavorite"`
}

func NaverSearch(query string, display int) (*NaverSearchResponse, error) {
    ClientID := os.Getenv("Naver_Client_ID")
    ClientSecret := os.Getenv("Naver_Secret")

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

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var searchResponse NaverSearchResponse
    err = json.Unmarshal(body, &searchResponse)
    if err != nil {
        return nil, err
    }
    
    return &searchResponse, nil
}

// controllers/searchController.go
func NaverSearchHandler(c *gin.Context) {
    query := c.Query("query")
    display := 5

    // 동시에 실행할 고루틴을 위한 에러 채널과 결과 채널 생성
    errChan := make(chan error, 2)
    searchChan := make(chan *NaverSearchResponse, 1)
    favoritesChan := make(chan []db.Restaurant, 1)

    // Naver API 검색을 고루틴으로 실행
    go func() {
        result, err := NaverSearch(query, display)
        if err != nil {
            errChan <- err
            return
        }
        searchChan <- result
    }()

    // 사용자의 즐겨찾기 목록을 동시에 가져오기
    session := sessions.Default(c)
    email := session.Get("email")
    go func() {
        if email != nil {
            favorites, err := db.GetRestaurantsDB(c.MustGet("mongoClient").(*mongo.Client), email.(string))
            if err != nil {
                errChan <- err
                return
            }
            favoritesChan <- favorites
        } else {
            favoritesChan <- []db.Restaurant{}
        }
    }()

    // 두 고루틴의 결과 대기
    var searchResult *NaverSearchResponse
    var userFavorites []db.Restaurant

    // 타임아웃 설정
    timeout := time.After(5 * time.Second)

    for i := 0; i < 2; i++ {
        select {
        case err := <-errChan:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        case searchResult = <-searchChan:
        case userFavorites = <-favoritesChan:
        case <-timeout:
            c.JSON(http.StatusRequestTimeout, gin.H{"error": "요청 시간이 초과되었습니다"})
            return
        }
    }

    // 결과 처리
    favoriteMap := make(map[string]bool)
    for _, fav := range userFavorites {
        favoriteMap[fav.Name] = true
    }

    restaurants := make([]RestaurantResponse, 0)
    for _, item := range searchResult.Items {
        name := strings.ReplaceAll(item.Title, "<b>", "")
        name = strings.ReplaceAll(name, "</b>", "")
        
        restaurants = append(restaurants, RestaurantResponse{
            Name:       name,
            Address:    item.RoadAddress,
            Category:   item.Category,
            IsFavorite: favoriteMap[name],
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "restaurants": restaurants,
    })
}