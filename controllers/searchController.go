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
    email := session.Get("email")
    var userFavorites []db.Restaurant
    if email != nil {
        userFavorites, _ = db.GetRestaurantsDB(c.MustGet("mongoClient").(*mongo.Client), email.(string))
    }

    // map은 반복문 밖에서 한 번만 생성
    favoriteMap := make(map[string]bool)
    for _, fav := range userFavorites {
        favoriteMap[fav.Name] = true
    }
    restaurants := make([]RestaurantResponse, 0)
    // 반복문에서는 생성된 map을 사용만 함
    for _, item := range result.Items {
        name := strings.ReplaceAll(item.Title, "<b>", "")
        name = strings.ReplaceAll(name, "</b>", "")
        
        isFavorite := favoriteMap[name]
        
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