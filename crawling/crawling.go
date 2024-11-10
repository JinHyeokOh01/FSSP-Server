package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
    
    "github.com/PuerkitoBio/goquery"
)

// Restaurant 구조체는 식당 정보를 저장합니다
type Restaurant struct {
    Name        string   `json:"name"`
    Menus       []Menu   `json:"menus"`
    OpeningHours string  `json:"opening_hours"`
    Photos      []string `json:"photos"`
}

// Menu 구조체는 메뉴 정보를 저장합니다
type Menu struct {
    Name  string `json:"name"`
    Price string `json:"price"`
}

func main() {
    // 예시 검색어와 지역으로 크롤링 실행
    restaurants, err := crawlNaverMap("강남 맛집", "서울 강남구")
    if err != nil {
        log.Fatal(err)
    }

    // 결과를 JSON으로 출력
    jsonData, err := json.MarshalIndent(restaurants, "", "    ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jsonData))
}

func crawlNaverMap(query, location string) ([]Restaurant, error) {
    // HTTP 클라이언트 설정
    client := &http.Client{
        Timeout: time.Second * 10,
    }

    // 네이버 지도 검색 URL 생성
    searchURL := fmt.Sprintf("https://map.naver.com/v5/search/%s%%20%s", query, location)

    // HTTP GET 요청
    req, err := http.NewRequest("GET", searchURL, nil)
    if err != nil {
        return nil, err
    }

    // User-Agent 설정
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // HTML 파싱
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }

    var restaurants []Restaurant

    // 식당 목록 순회
    doc.Find(".list_item").Each(func(i int, s *goquery.Selection) {
        restaurant := Restaurant{}
        
        // 식당 이름 추출
        restaurant.Name = strings.TrimSpace(s.Find(".name").Text())
        
        // 영업시간 추출
        restaurant.OpeningHours = strings.TrimSpace(s.Find(".time").Text())
        
        // 메뉴 정보 추출
        s.Find(".menu_list li").Each(func(i int, s *goquery.Selection) {
            menu := Menu{
                Name:  strings.TrimSpace(s.Find(".menu_name").Text()),
                Price: strings.TrimSpace(s.Find(".menu_price").Text()),
            }
            restaurant.Menus = append(restaurant.Menus, menu)
        })
        
        // 사진 URL 추출
        s.Find(".thumb img").Each(func(i int, s *goquery.Selection) {
            if photoURL, exists := s.Attr("src"); exists {
                restaurant.Photos = append(restaurant.Photos, photoURL)
            }
        })
        
        restaurants = append(restaurants, restaurant)
    })

    return restaurants, nil
}

// 에러 처리를 위한 사용자 정의 에러 타입
type CrawlError struct {
    Message string
    Err     error
}

func (e *CrawlError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}