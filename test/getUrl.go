package main

import (
    "fmt"
    "log"
    "time"

    "github.com/go-rod/rod"
    "github.com/go-rod/rod/lib/launcher"
)

func getNaverMapURL(searchKeyword string) (string, error) {
    // 브라우저 설정
    u := launcher.New().
        Leakless(false).
        Headless(true).  // 브라우저 창을 보이지 않게 설정
        MustLaunch()

    browser := rod.New().
        ControlURL(u).
        MustConnect()
    defer browser.MustClose()

    page := browser.MustPage("https://map.naver.com/p")

    // 페이지 로딩 대기
    page.MustWaitLoad()

    // 검색 입력
    searchInput := page.MustElement("#search-input")
    searchInput.MustInput(searchKeyword)
    
    // 검색 버튼 클릭
    page.MustElement("button.btn_search").MustClick()

    // 최소한의 대기 시간
    time.Sleep(time.Millisecond * 800)

    // 첫 번째 결과 클릭
    page.MustElement(".lst_site .item_info").MustClick()

    // 최소한의 대기 시간
    time.Sleep(time.Millisecond * 500)

    return page.MustInfo().URL, nil
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