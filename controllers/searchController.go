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

func NaverSearchHandler(c *gin.Context) {
    query := c.Query("query")
    display := 5

    // 1. Naver API 검색 결과 가져오기
    result, err := NaverSearch(query, display)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 2. 사용자의 즐겨찾기 목록 가져오기
    session := sessions.Default(c)
    email := session.Get("userEmail")
    var userFavorites []db.Restaurant
    if email != nil {
        userFavorites, _ = db.GetRestaurantsDB(c.MustGet("mongoClient").(*mongo.Client), email.(string))
    }

    // 3. 응답 포맷 변환
    restaurants := make([]RestaurantResponse, 0)
    for _, item := range result.Items {
        // HTML 태그 제거 및 특수문자 처리
        name := strings.ReplaceAll(item.Title, "<b>", "")
        name = strings.ReplaceAll(name, "</b>", "")

        // 즐겨찾기 여부 확인
        isFavorite := false
        for _, fav := range userFavorites {
            if fav.Name == name {
                isFavorite = true
                break
            }
        }

        restaurants = append(restaurants, RestaurantResponse{
            Name:       name,
            Address:    item.RoadAddress,
            Category:   item.Category,
            IsFavorite: isFavorite,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "restaurants": restaurants,
    })
}