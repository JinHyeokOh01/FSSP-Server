//Claude
/*
package main

import (
    "fmt"
    "log"
    "github.com/chromedp/chromedp"
    "context"
    "time"
)
/*
func main() {
    // 크롬 옵션 설정
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
    )
    
    // 컨텍스트 생성
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()
    
    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()
    
    // 타임아웃 설정
    ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
    defer cancel()

    // URL을 저장할 변수
    var href string
    
    // 웹페이지 방문 및 URL 추출
    err := chromedp.Run(ctx,
        // 웹페이지 방문 (실제 웹사이트 URL로 변경 필요)
        chromedp.Navigate("https://map.naver.com/p/search/%EA%B7%B8%EB%A1%9C%EB%98%90"),
        
        // CSS 선택자로 요소 찾기 및 href 속성 가져오기
        chromedp.AttributeValue("#_pcmap_list_scroll_container > ul > li:nth-child(1) > div.qbGlu > div.ouxiq > a:nth-child(2)", "href", &href, nil),
    )
    
    if err != nil {
        log.Fatal(err)
    }
    
    // 결과 출력
    fmt.Printf("찾은 URL: %s\n", href)
}
    */
/*
func main() {
    // Chrome 실행을 위한 컨텍스트 생성
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    // 타임아웃 설정
    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // 크롤링할 URL
    url := "https://map.naver.com/p/search/%EA%B7%B8%EB%A1%9C%EB%98%90"
    var newURL string
    clickXPath := "//*[@id=\"_pcmap_list_scroll_container\"]/ul/li[1]"

    // 작업 실행
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.WaitVisible(clickXPath, chromedp.ByQuery),
        chromedp.Click(clickXPath, chromedp.NodeVisible),
        chromedp.Sleep(3*time.Second), // 클릭 후 로드 대기
        chromedp.Location(&newURL),     // 현재 URL 가져오기
    )
    if err != nil {
        log.Fatal(err)
    }

    // 결과 출력
    fmt.Println("이동한 URL:", newURL)
}
	*/
/*
func main() {
    // Chrome 옵션 설정 - User-Agent 추가 및 기타 설정
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
        chromedp.Flag("disable-gpu", true),
        chromedp.Flag("no-sandbox", true),
        chromedp.Flag("disable-dev-shm-usage", true),
        // User-Agent 설정
        chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
        // 추가 플래그들
        chromedp.Flag("disable-web-security", true),
        chromedp.Flag("disable-background-networking", true),
        chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
        chromedp.Flag("disable-background-timer-throttling", true),
        chromedp.Flag("disable-backgrounding-occluded-windows", true),
        chromedp.Flag("disable-breakpad", true),
        chromedp.Flag("disable-client-side-phishing-detection", true),
        chromedp.Flag("disable-default-apps", true),
        chromedp.Flag("disable-extensions", true),
        chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
        chromedp.Flag("disable-hang-monitor", true),
        chromedp.Flag("disable-ipc-flooding-protection", true),
        chromedp.Flag("disable-popup-blocking", true),
        chromedp.Flag("disable-prompt-on-repost", true),
        chromedp.Flag("disable-renderer-backgrounding", true),
        chromedp.Flag("disable-sync", true),
        chromedp.Flag("force-color-profile", "srgb"),
        chromedp.Flag("metrics-recording-only", true),
        chromedp.Flag("no-first-run", true),
    )

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    // Chrome 컨텍스트 생성 (로깅 추가)
    ctx, cancel := chromedp.NewContext(
        allocCtx,
        chromedp.WithLogf(log.Printf),
    )
    defer cancel()

    // 타임아웃 증가
    ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
    defer cancel()

    url := "https://pcmap.place.naver.com/restaurant/20506331/home"
    var texts []string

    err := chromedp.Run(ctx,
        // 페이지 이동 전 쿠키 및 캐시 초기화
        chromedp.EmulateViewport(1920, 1080),
        chromedp.Navigate(url),
        
        // 로딩 대기 시간 증가
        //chromedp.Sleep(1*time.Second),
        
        // DOM이 완전히 로드될 때까지 대기
        chromedp.WaitReady("body", chromedp.ByQuery),
        
        // 특정 요소가 보일 때까지 대기
        chromedp.WaitVisible(".O8qbU.pSavy div a", chromedp.ByQuery),
        
        // 클릭
        chromedp.Click(".O8qbU.pSavy div a", chromedp.ByQuery),
        
        // 클릭 후 대기 시간 증가
        //chromedp.Sleep(1*time.Second),
        
        // JavaScript로 텍스트 추출
        chromedp.Evaluate(`
            Array.from(document.querySelectorAll('.A_cdD')).map(el => el.textContent)
        `, &texts),
    )

    if err != nil {
        log.Printf("크롤링 중 에러 발생: %v\n", err)
        return
    }

    fmt.Println("크롤링된 식당 정보:")
    for i, text := range texts {
        fmt.Printf("%d. %s\n", i+1, text)
    }
}
*/