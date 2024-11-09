/*

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/antchfx/htmlquery"
)

func main() {
	// 크롤러 생성
	c := colly.NewCollector()

	// 크롤링할 iframe의 ID를 지정
	iframeID := "entryIframe" // 원하는 iframe ID로 변경하세요.

	// 특정 URL에 접근
	c.OnHTML("iframe", func(e *colly.HTMLElement) {
		// iframe의 ID가 일치하는지 확인
		if e.Attr("id") == iframeID {
			// iframe의 src 속성 가져오기
			iframeSrc := e.Attr("src")
			if !strings.HasPrefix(iframeSrc, "http") {
				iframeSrc = e.Request.AbsoluteURL(iframeSrc)
			}

			// iframe 페이지 크롤링
			e.Request.Visit(iframeSrc)
		}
	})

	// iframe 내부 페이지에서 특정 클래스 크롤링
	c.OnHTML("body", func(e *colly.HTMLElement) {
		htmlContent, err := e.DOM.Html() // HTML 내용 가져오기
		if err != nil {
			log.Fatal(err) // 오류 처리
		}

		doc, err := htmlquery.Parse(strings.NewReader(htmlContent)) // HTML 파싱
		if err != nil {
			log.Fatal(err) // 오류 처리
		}

		// XPath로 특정 클래스 선택 (예: .my-class)
		xpathExpr := "//div[contains(@class, 'my-class')]"
		nodes := htmlquery.Find(doc, xpathExpr)

		for _, node := range nodes {
			fmt.Println(htmlquery.OutputHTML(node, true))
		}
	})

	// 크롤링 시작
	err := c.Visit("https://map.naver.com/p/search/%EA%B7%B8%EB%A1%9C%EB%98%90/place/20506331") // 크롤링할 페이지 URL
	if err != nil {
		log.Fatal(err)
	}
}
*/

//이거 잘 가지고 있어야 함
/*
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// 새로운 컨텍스트 생성
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 타임아웃 설정
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var xpath string
	var elementText string

	xpath = "//*[@id=\"app-root\"]"

	// 크롤링할 URL을 설정합니다.
	url := "https://pcmap.place.naver.com/restaurant/20506331/home" // 원하는 웹 페이지로 변경하세요.

	// 작업 실행
	err := chromedp.Run(ctx,
		chromedp.Navigate(url), // 페이지 탐색
		//chromedp.WaitVisible(xpath, chromedp.ByXPath), // 요소가 보일 때까지 대기
		chromedp.Text(xpath, &elementText, chromedp.NodeVisible), // 텍스트 가져오기
	)
	
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("XPath '%s'의 텍스트: %s\n", xpath, elementText)
}

*/

/*
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/chromedp/chromedp"
	//"github.com/chromedp/cdproto/cdp"
)

func main() {
    // Chrome 실행을 위한 컨텍스트 생성
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    // 타임아웃 설정
    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // 크롤링할 URL
    url := "https://pcmap.place.naver.com/restaurant/20506331/home"
    var elementText string
    clickXPath := "//*[@id=\"app-root\"]/div/div/div/div[6]/div/div[2]/div[1]/div/div[3]/div/a/div[1]"
    textXPath := "//*[@id=\"app-root\"]/div/div/div/div[6]/div/div[2]/div[1]/div/div[3]/div/a/div[2]/div/span[1]/div"

    // 작업 실행
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.WaitReady(clickXPath, chromedp.ByXPath), // 요소가 로드될 때까지 대기
        chromedp.Click(clickXPath, chromedp.NodeVisible),
        chromedp.Sleep(3*time.Second), // 클릭 후 로드 대기
        chromedp.WaitVisible(textXPath, chromedp.ByXPath), // 텍스트 요소가 보일 때까지 대기
        chromedp.Text(textXPath, &elementText, chromedp.NodeVisible),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 결과 출력
    fmt.Println("가져온 텍스트:", elementText)
}
*/


//CSS Selector
/*
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/chromedp/chromedp"
)

func main() {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    url := "https://pcmap.place.naver.com/restaurant/20506331/home"
    var elementText string

    // 클릭할 요소의 CSS 선택자
    clickSelector := "#app-root > div > div > div > div:nth-child(5) > div > div:nth-child(2) > div.place_section_content > div > div.O8qbU.pSavy > div > a > div"
    // 클릭 후 보이는 요소의 CSS 선택자
    textSelector := "#app-root > div > div > div > div:nth-child(5) > div > div:nth-child(2) > div.place_section_content > div > div.O8qbU.pSavy > div > a > div:nth-child(2) > div > span.A_cdD" // 여기에 클릭 후 보이는 요소의 CSS 선택자를 입력하세요.

    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.WaitVisible(clickSelector, chromedp.ByQuery), // 요소가 보일 때까지 대기
        chromedp.Click(clickSelector, chromedp.NodeVisible), // 클릭
        chromedp.Sleep(2*time.Second), // 클릭 후 로드 대기
        chromedp.WaitVisible(textSelector, chromedp.ByQuery), // 보일 요소가 보일 때까지 대기
        chromedp.Text(textSelector, &elementText, chromedp.NodeVisible), // 텍스트 가져오기
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("가져온 텍스트:", elementText)
}
*/

//Claude
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/chromedp/chromedp"
)

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
        chromedp.Sleep(5*time.Second),
        
        // DOM이 완전히 로드될 때까지 대기
        chromedp.WaitReady("body", chromedp.ByQuery),
        
        // 특정 요소가 보일 때까지 대기
        chromedp.WaitVisible(".O8qbU.pSavy div a", chromedp.ByQuery),
        
        // 클릭
        chromedp.Click(".O8qbU.pSavy div a", chromedp.ByQuery),
        
        // 클릭 후 대기 시간 증가
        chromedp.Sleep(5*time.Second),
        
        // JavaScript로 텍스트 추출
        chromedp.Evaluate(`
            Array.from(document.querySelectorAll('.A_cdD')).map(el => el.textContent)
        `, &texts),
    )

    if err != nil {
        log.Printf("크롤링 중 에러 발생: %v\n", err)
        return
    }

    fmt.Println("크롤링된 메뉴 정보:")
    for i, text := range texts {
        fmt.Printf("%d. %s\n", i+1, text)
    }
}