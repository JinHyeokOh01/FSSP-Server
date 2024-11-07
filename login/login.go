package login

import(
	"net/http"
    "time"
    "math/rand"
    "encoding/base64"
	"encoding/json"
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
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
	}
}

func GoogleForm(c *gin.Context) {
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
 
 func GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
 
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
 }
 
 func GoogleLoginHandler(c *gin.Context) {
 
	state := GenerateStateOauthCookie(c.Writer)
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect,url)
 }

 func GoogleAuthCallback(c *gin.Context) {
	oauthstate, _ := c.Request.Cookie("oauthstate")
 
	if c.Request.FormValue("state") != oauthstate.Value {
	   log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
	   c.Redirect(http.StatusTemporaryRedirect,"/")
	   return
	}
 
	data, err := GetGoogleUserInfo(c.Request.FormValue("code"))
	if err != nil {
	   log.Println(err.Error())
	   c.Redirect( http.StatusTemporaryRedirect,"/")
	   return
	}
 
	fmt.Fprint(c.Writer, string(data))
 }
 
 func GetGoogleUserInfo(code string) ([]byte, error) {
	const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" 
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
	   return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}
 
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {                                             
	   return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	src_json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        return nil, fmt.Errorf("Failed to unmarshal JSON:", err.Error())
    }
	defer resp.Body.Close()

	var n map[string]interface{}
	err = json.Unmarshal(src_json, &n)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(n["given_name"])

	return src_json, err
 }
/*
 func main() {

	r := gin.Default()
 
	r.GET("/", GoogleForm)
	r.GET("/auth/google/login", GoogleLoginHandler)
	r.GET("/auth/google/callback", GoogleAuthCallback)
 
	r.Run(":5000")
 }
*/