package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/chromedp/chromedp"
)

func getNaverMapURL(searchKeyword string) (string, error) {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
    defer cancel()

    var currentURL string
    err := chromedp.Run(ctx,
        // 네이버 지도로 이동
        chromedp.Navigate("https://map.naver.com/p/search/그로또"),
        
        // 검색창이 나타날 때까지 대기
        chromedp.WaitVisible("input.input_search", chromedp.ByQuery),
        
        // 검색어 입력 (class로 선택)
        chromedp.SendKeys("input.input_search", searchKeyword, chromedp.ByQuery),
        
        // 검색 버튼 클릭 (class로 선택)
        chromedp.Click("button.button_search", chromedp.ByQuery),
        
        // 결과 로딩 대기
        chromedp.Sleep(time.Millisecond*800),
        
        // 첫 번째 결과 클릭
        chromedp.Click(".ouxiq div a", chromedp.ByQuery),
        
        // URL 변경 대기
        chromedp.Sleep(time.Millisecond*500),
        
        // 현재 URL 가져오기
        chromedp.Location(&currentURL),
    )

    if err != nil {
        return "", err
    }

    return currentURL, nil
}

func main() {
    fmt.Print("검색할 장소: ")
    var keyword string
    fmt.Scanln(&keyword)

    url, err := getNaverMapURL(keyword)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("URL:", url)
}