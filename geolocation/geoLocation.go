/*
package gl

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"

	"github.com/joho/godotenv"

)

// GeolocationResponse는 NCP Geolocation API의 응답 구조체입니다
type GeolocationResponse struct {
    ReturnCode    int    `json:"returnCode"`
    ReturnMessage string `json:"returnMessage"`
    GeoLocation   struct {
        Country    string  `json:"country"`
        Code      string  `json:"r1"`
        Region    string  `json:"r2"`
        City      string  `json:"r3"`
        Latitude  float64 `json:"lat"`
        Longitude float64 `json:"long"`
        Net       string  `json:"net"`
        IPAddress string  `json:"ip"`
    } `json:"geoLocation"`
}

// getPublicIP는 외부 IP 주소를 가져옵니다
func getPublicIP() (string, error) {
    services := []string{
        "https://api.ipify.org",
        "https://ifconfig.me/ip",
        "https://api.ipify.org?format=text",
    }
    
    var lastErr error
    for _, service := range services {
        resp, err := http.Get(service)
        if err != nil {
            lastErr = err
            continue
        }
        defer resp.Body.Close()
        
        ip, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            lastErr = err
            continue
        }
        
        publicIP := strings.TrimSpace(string(ip))
        if publicIP != "" {
            return publicIP, nil
        }
    }
    
    return "", fmt.Errorf("public IP 조회 실패: %v", lastErr)
}

// makeSignature는 API Gateway 서명을 생성합니다
func makeSignature(timestamp, accessKey, secretKey string) string {
    method := "GET"
    uri := "/geolocation/v2/geoLocation"
    
    message := fmt.Sprintf("%s %s\n%s\n%s", method, uri, timestamp, accessKey)
    
    mac := hmac.New(sha256.New, []byte(secretKey))
    mac.Write([]byte(message))
    
    signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
    return signature
}

func getGeolocation(accessKey, secretKey, ipAddress string) (*GeolocationResponse, error) {
    baseURL := "https://geolocation.apigw.ntruss.com/geolocation/v2/geoLocation"
    params := url.Values{}
    params.Add("ip", ipAddress)
    params.Add("ext", "t")
    params.Add("enc", "utf8")
    params.Add("responseFormatType", "json")
    
    fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

    timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
    signature := makeSignature(timestamp, accessKey, secretKey)

    client := &http.Client{}
    req, err := http.NewRequest("GET", fullURL, nil)
    if err != nil {
        return nil, fmt.Errorf("요청 생성 실패: %v", err)
    }

    req.Header.Set("x-ncp-apigw-timestamp", timestamp)
    req.Header.Set("x-ncp-iam-access-key", accessKey)
    req.Header.Set("x-ncp-apigw-signature-v2", signature)

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("API 요청 실패: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("응답 읽기 실패: %v", err)
    }

    var result GeolocationResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("JSON 파싱 실패: %v", err)
    }

    return &result, nil
}

func main() {
    // 환경 변수에서 인증 정보 가져오기
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    ClientID := os.Getenv("Naver_Client_ID")
    ClientSecret := os.Getenv("Naver_Secret")

    // 현재 컴퓨터의 공인 IP 주소 가져오기
    publicIP, err := getPublicIP()
    if err != nil {
        log.Fatalf("공인 IP 조회 실패: %v\n", err)
    }
    fmt.Printf("현재 공인 IP 주소: %s\n\n", publicIP)

    // 위치 정보 조회
    result, err := getGeolocation(ClientID, ClientSecret, publicIP)
    if err != nil {
        log.Fatalf("위치 정보 조회 실패: %v\n", err)
    }

    // 결과 출력
    if result.ReturnCode == 0 {
        fmt.Println("위치 정보:")
        fmt.Printf("국가: %s\n", result.GeoLocation.Country)
        fmt.Printf("지역: %s\n", result.GeoLocation.Region)
        fmt.Printf("도시: %s\n", result.GeoLocation.City)
        fmt.Printf("위도: %f\n", result.GeoLocation.Latitude)
        fmt.Printf("경도: %f\n", result.GeoLocation.Longitude)
        fmt.Printf("네트워크: %s\n", result.GeoLocation.Net)
    } else {
        fmt.Printf("API 오류: %s (코드: %d)\n", result.ReturnMessage, result.ReturnCode)
    }
}
	*/

	package geolocation

	import (
		"crypto/hmac"
		"crypto/sha256"
		"encoding/base64"
		"encoding/json"
		"fmt"
		//"io/ioutil"
		"log"
		"net/http"
		//"net/url"
		"strconv"
		"time"
		//"os"

		"github.com/joho/godotenv"
        "github.com/gin-gonic/gin"
	)
	
	// GeolocationResponse는 NCP Geolocation API의 응답 구조체입니다
	type GeolocationResponse struct {
		ReturnCode    int    `json:"returnCode"`
		ReturnMessage string `json:"returnMessage"`
		GeoLocation   struct {
			Country    string  `json:"country"`
			Code      string  `json:"r1"`
			Region    string  `json:"r2"`
			City      string  `json:"r3"`
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"long"`
			Net       string  `json:"net"`
			IPAddress string  `json:"ip"`
		} `json:"geoLocation"`
	}
	
// makeSignature는 HMAC 서명을 생성합니다.
func makeSignature(method, basestring, timestamp, accessKey, secretKey string) string {
	// 메시지 생성
	message := fmt.Sprintf("%s %s\n%s\n%s", method, basestring, timestamp, accessKey)

	// HMAC 생성
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))

	// 서명 생성 및 Base64 인코딩
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

// requestAPI는 API 요청을 수행합니다.
// requestAPI는 API 요청을 수행하고 결과를 반환합니다.
func requestAPI(timestamp, accessKey, signature, uri string) (*GeolocationResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %v", err)
	}

	// 헤더 설정
	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", accessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", signature)

	// API 요청
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 상태 확인
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 요청 실패, 상태 코드: %d", resp.StatusCode)
	}

	// 응답 본문 읽기
	var result GeolocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("응답 JSON 파싱 실패: %v", err)
	}

	return &result, nil
}

func getGeolocation() (*GeolocationResponse, error) {
	// Signature 생성에 필요한 항목
	method := "GET"
	basestring := "/geolocation/v2/geoLocation?ip=218.209.223.78&ext=t&responseFormatType=json"
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	accessKey := "ncp_iam_BPAMKR4F7WtKdTFCLr2f" // access key id (from portal or sub account)
	secretKey := "ncp_iam_BPKMKRLi5RVaBIBWFT0jDKVWmWL8bLc0p5" // secret key (from portal or sub account)

	// 서명 생성
	signature := makeSignature(method, basestring, timestamp, accessKey, secretKey)

	// 결과 출력 (디버깅용)
	fmt.Println("타임스탬프:", timestamp)
	fmt.Println("액세스 키:", accessKey)
	fmt.Println("서명:", signature)

	// GET 요청
	hostname := "https://geolocation.apigw.ntruss.com"
	requestUri := hostname + basestring
	response, err := requestAPI(timestamp, accessKey, signature, requestUri)
	if err != nil {
		return nil, err
	}

	return response, nil
}
	
	func GeoLocationHandler(c *gin.Context) {
		// 인증 정보 직접 입력

		err := godotenv.Load()
    	if err != nil {
       		log.Fatal("Error loading .env file")
    	}

		//accessKey := os.Getenv("Naver_Cloud_Access_Key") // 여기에 실제 Access Key를 입력하세요
		//secretKey := os.Getenv("Naver_Cloud_Secret_Key") // 여기에 실제 Secret Key를 입력하세요
	
		// 테스트용 IP 주소
		ipAddress := "218.209.223.78"  // 테스트할 IP 주소
	
		fmt.Printf("테스트할 IP 주소: %s\n", ipAddress)
	
		// 위치 정보 조회
		result, err := getGeolocation()
		if err != nil {
			log.Fatalf("위치 정보 조회 실패: %v\n", err)
		}

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        c.JSON(http.StatusOK, result)
	
		// 결과 출력
		if result.ReturnCode == 0 {
			fmt.Println("\n[위치 정보]")
			fmt.Printf("국가: %s\n", result.GeoLocation.Country)
			fmt.Printf("지역: %s\n", result.GeoLocation.Region)
			fmt.Printf("도시: %s\n", result.GeoLocation.City)
			fmt.Printf("위도: %f\n", result.GeoLocation.Latitude)
			fmt.Printf("경도: %f\n", result.GeoLocation.Longitude)
			fmt.Printf("네트워크: %s\n", result.GeoLocation.Net)
		} else {
			fmt.Printf("\n[API 오류]\n메시지: %s\n코드: %d\n", result.ReturnMessage, result.ReturnCode)
		}
	}