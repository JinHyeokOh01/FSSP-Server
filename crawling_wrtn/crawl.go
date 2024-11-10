package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func main() {
	// Selenium WebDriver 설정
	const (
		seleniumPath = "path/to/selenium-server-standalone.jar" // Selenium 서버 경로
		chromeDriverPath = "path/to/chromedriver" // ChromeDriver 경로
		port = 8080
	)

	// Selenium 서버 시작
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.Output(nil), // Log output
	}
	srv, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Fatalf("Error starting the Selenium server: %v", err)
	}
	defer srv.Stop()

	// WebDriver 설정
	caps := selenium.Capabilities{"browserName": "chrome"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver: %v", err)
	}
	defer driver.Quit()

	// 네이버 지도 페이지 열기
	if err := driver.Get("https://map.naver.com/p/search/그로또"); err != nil {
		log.Fatalf("Failed to load page: %v", err)
	}

	// 페이지 로딩 대기
	time.Sleep(5 * time.Second) // 충분한 시간 대기

	// 식당 정보 추출 (CSS 선택자는 실제 페이지 구조에 맞게 수정 필요)
	restaurants := driver.FindElements(selenium.ByCSSSelector(".place_section"))
	for _, restaurant := range restaurants {
		name, _ := restaurant.FindElement(selenium.ByCSSSelector(".place_name"))
		menu, _ := restaurant.FindElement(selenium.ByCSSSelector(".menu_item"))
		openHours, _ := restaurant.FindElement(selenium.ByCSSSelector(".open_hours"))

		nameText, _ := name.Text()
		menuText, _ := menu.Text()
		openHoursText, _ := openHours.Text()

		fmt.Printf("식당 이름: %s\n", nameText)
		fmt.Printf("메뉴: %s\n", menuText)
		fmt.Printf("영업시간: %s\n", openHoursText)
		fmt.Println("=================================")
	}
}
