//Claude
package crwl

import (
    "fmt"
    "log"
    "github.com/chromedp/chromedp"
    "github.com/gin-gonic/gin"

    "context"
    "time"
    "net/http"
)

/*
type TextElement struct {
    Index int    `json:"index"`
    Text  string `json:"text"`
}
    */

func crawling() ([]string, error) {
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
        chromedp.Sleep(1*time.Second),

        chromedp.WaitVisible(".ipNNM .ZHqBk img.K0PDV"),
        
        // JavaScript로 텍스트 추출
        chromedp.Evaluate(`
                (() => {
            const texts = Array.from(document.querySelectorAll('.A_cdD'))
                .map(el => el.textContent.trim());
                
            const numbers = Array.from(document.querySelectorAll('.xlx7Q'))
                .map(el => el.textContent.trim());

            const menu = Array.from(document.querySelectorAll('.MN48z'))
                .map(el => el.textContent.trim());

            const imgSrc = Array.from(document.querySelectorAll('.place_section_content img'))
                .map(el => el.getAttribute('src'))
            return [...texts, ...numbers, ...menu, ...imgSrc];  // spread 연산자를 사용하여 두 배열을 올바르게 합침
        })()
    `, &texts),
    )

    if err != nil {
        log.Printf("크롤링 중 에러 발생: %v\n", err)
        return []string{}, err
    }

    return texts, nil
}

func CrawlingHandler(c *gin.Context) {
    // 크롤링 조회
    result, err := crawling()
    if err != nil {
        log.Fatalf("정보 조회 실패: %v\n", err)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)

    fmt.Println(result)
}