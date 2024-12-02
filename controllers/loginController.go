package controllers

import (
    "context"
    "encoding/base64"
    "encoding/json"
	"encoding/gob"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "os"
    "time"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/JinHyeokOh01/FSSP-Server/db"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

const (
    userKey = "user"
)

var googleOauthConfig oauth2.Config

func InitAuth(router *gin.Engine) error {
    // gob 등록 추가
    gob.Register(map[string]interface{}{})
    
    // 세션 미들웨어 설정
    sessionSecret := os.Getenv("SESSION_SECRET")
    if sessionSecret == "" {
        return fmt.Errorf("SESSION_SECRET environment variable is not set")
    }
    
    store := cookie.NewStore([]byte(sessionSecret))
    store.Options(sessions.Options{
        MaxAge: 86400 * 7, // 7일
        Path:   "/",
        Secure: false, // HTTPS를 사용하는 경우 true로 설정
        HttpOnly: true,
    })
    router.Use(sessions.Sessions("mysession", store))
    
    // OAuth 설정
    clientID := os.Getenv("Client_ID")
    clientSecret := os.Getenv("Client_Secret")
    if clientID == "" || clientSecret == "" {
        return fmt.Errorf("Client_ID or Client_Secret environment variable is not set")
    }

    googleOauthConfig = oauth2.Config{
        RedirectURL:  "http://localhost:5000/auth/google/callback",
        ClientID:     clientID,
        ClientSecret: clientSecret,
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.profile",
            "https://www.googleapis.com/auth/userinfo.email",
        },
        Endpoint: google.Endpoint,
    }
    
    return nil
}

func GoogleForm(c *gin.Context) {
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <html>
            <head>
                <title>Go Oauth2.0 Test</title>
            </head>
            <body>
                <p><a href='./auth/google/login'>Google Login</a></p>
            </body>
        </html>
    `))
}

func GenerateStateOauthCookie(w http.ResponseWriter) string {
    expiration := time.Now().Add(24 * time.Hour)
    b := make([]byte, 16)
    rand.Read(b)
    state := base64.URLEncoding.EncodeToString(b)
    cookie := &http.Cookie{
        Name:     "oauthstate",
        Value:    state,
        Expires:  expiration,
        HttpOnly: true,
        Path:     "/",
    }
    http.SetCookie(w, cookie)
    return state
}

func GoogleLoginHandler(c *gin.Context) {
    state := GenerateStateOauthCookie(c.Writer)
    url := googleOauthConfig.AuthCodeURL(state)
    c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleAuthCallback(c *gin.Context, client *mongo.Client) {
    oauthstate, err := c.Request.Cookie("oauthstate")
    if err != nil {
        handleError(c, http.StatusBadRequest, "Missing oauth state cookie", err)
        return
    }

    if c.Request.FormValue("state") != oauthstate.Value {
        handleError(c, http.StatusBadRequest, "Invalid oauth state", nil)
        return
    }

    data, err := GetGoogleUserInfo(c.Request.FormValue("code"))
    if err != nil {
        handleError(c, http.StatusUnauthorized, "Failed to get user info", err)
        return
    }

    var userData map[string]interface{}
    if err := json.Unmarshal(data, &userData); err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to parse user data", err)
        return
    }

    // MongoDB 사용자 처리
    err = handleMongoDBUser(client, userData)
    if err != nil {
        handleError(c, http.StatusInternalServerError, "Database error", err)
        return
    }

    // 세션 저장
    session := sessions.Default(c)
    session.Set(userKey, userData)
    if err := session.Save(); err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to save session", err)
        return
    }

    // 로그인 성공 후 메인 페이지로 리다이렉트
    c.Redirect(http.StatusTemporaryRedirect, "/")
}

func GetGoogleUserInfo(code string) ([]byte, error) {
    const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
    
    token, err := googleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        return nil, fmt.Errorf("failed to exchange token: %v", err)
    }

    resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %v", err)
    }
    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}

func LogoutHandler(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Options(sessions.Options{
        MaxAge: -1,
        Path:   "/",
    })
    if err := session.Save(); err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to clear session", err)
        return
    }
    c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := sessions.Default(c)
        user := session.Get(userKey)
        if user == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        c.Set(userKey, user)
        c.Next()
    }
}

func GetCurrentUser(c *gin.Context) map[string]interface{} {
    session := sessions.Default(c)
    user := session.Get(userKey)
    if user == nil {
        return nil
    }
    return user.(map[string]interface{})
}

// 에러 처리 헬퍼 함수
func handleError(c *gin.Context, status int, message string, err error) {
    if err != nil {
        log.Printf("%s: %v", message, err)
    }
    c.JSON(status, gin.H{"error": message})
}

// MongoDB 사용자 처리 헬퍼 함수
func handleMongoDBUser(client *mongo.Client, userData map[string]interface{}) error {
    var result bson.M
    coll := client.Database("FSSP_DB").Collection("users")
    filter := bson.M{"email": userData["email"].(string)}
    
    err := coll.FindOne(context.TODO(), filter).Decode(&result)
    if err == mongo.ErrNoDocuments {
        return db.UserAddDB(client, userData["email"].(string), userData["given_name"].(string))
    }
    return err
}