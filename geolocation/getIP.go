package geolocation

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)
/*
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
	*/

// getPublicIP는 외부 IP 주소를 가져옵니다
func GetPublicIP() (string, error) {
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