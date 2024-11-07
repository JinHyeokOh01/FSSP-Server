package lgn

import(
	"net/http"
    "time"
    "math/rand"
    "encoding/base64"
    "log"
	"fmt"
    "context"
    "io/ioutil"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/joho/godotenv"
	"os"
)

var (
    googleOauthConfig oauth2.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:5000/auth/google/callback",
	ClientID:     os.Getenv("Client_ID"),
	ClientSecret: os.Getenv("Client_Secret"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
	}
}

func googleForm(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
	   "<html>" +
		  "\n<head>\n    " +
		  "<title>Go Oauth2.0 Test</title>\n" +
		  "</head>\n" +
		  "<body>\n<p>" +
		  "<a href='./auth/google/login'>Google Login</a>" +
		  "</p>\n" +
		  "</body>\n" +
	   "</html>"))
 }
 
 func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
 
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
 }
 
 func googleLoginHandler(c *gin.Context) {
 
	state := generateStateOauthCookie(c.Writer)
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect,url)
 }

 func googleAuthCallback(c *gin.Context) {
	oauthstate, _ := c.Request.Cookie("oauthstate") // 12
 
	if c.Request.FormValue("state") != oauthstate.Value { // 13
	   log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
	   c.Redirect(http.StatusTemporaryRedirect,"/")
	   return
	}
 
	data, err := getGoogleUserInfo(c.Request.FormValue("code")) // 14
	if err != nil {                                     // 15
	   log.Println(err.Error())
	   c.Redirect( http.StatusTemporaryRedirect,"/")
	   return
	}
 
	fmt.Fprint(c.Writer, string(data)) // 16
 }
 
 func getGoogleUserInfo(code string) ([]byte, error) { // 17
	const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" 
	token, err := googleOauthConfig.Exchange(context.Background(), code) // 18
	if err != nil {                                                      // 19
	   return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}
 
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken) // 20
	if err != nil {                                              // 21
	   return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}
 
	return ioutil.ReadAll(resp.Body) // 23
 }