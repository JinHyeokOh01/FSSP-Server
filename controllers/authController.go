// controllers/auth_controller.go
package controllers

import (
    "context"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
    db *mongo.Database
}

func NewAuthController(db *mongo.Database) *AuthController {
    return &AuthController{db: db}
}

func (ac *AuthController) isEmailExists(email string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    count, err := ac.db.Collection("users").CountDocuments(ctx, bson.M{"email": email})
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (ac *AuthController) Register(c *gin.Context) {
    var input struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
        Name     string `json:"name" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 입력입니다"})
        return
    }

    // 이메일 중복 체크
    exists, err := ac.isEmailExists(input.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류가 발생했습니다"})
        return
    }
    if exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "이미 등록된 이메일입니다"})
        return
    }

    // 비밀번호 해싱
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 처리 중 오류가 발생했습니다"})
        return
    }

    now := time.Now()
    user := bson.M{
        "email":       input.Email,
        "password":    string(hashedPassword),
        "name":        input.Name,
        "createdAt":   now,
        "updatedAt":   now,
        "restaurants": []interface{}{},
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = ac.db.Collection("users").InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "사용자 생성에 실패했습니다"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "회원가입이 완료되었습니다"})
}

func (ac *AuthController) Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 입력입니다"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user bson.M
    err := ac.db.Collection("users").FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "이메일 또는 비밀번호가 올바르지 않습니다"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류가 발생했습니다"})
        return
    }

    // 비밀번호 검증
    err = bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(input.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "이메일 또는 비밀번호가 올바르지 않습니다"})
        return
    }

    // 세션에 사용자 정보 저장
    session := sessions.Default(c)
    session.Set("userId", user["_id"].(primitive.ObjectID).Hex())
    session.Set("email", user["email"])
    session.Set("name", user["name"])
    if err := session.Save(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "세션 저장에 실패했습니다"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "로그인 성공",
        "user": gin.H{
            "id":    user["_id"].(primitive.ObjectID).Hex(),
            "email": user["email"],
            "name":  user["name"],
        },
    })
}

func (ac *AuthController) Logout(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Options(sessions.Options{MaxAge: -1})
    if err := session.Save(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "로그아웃 처리 중 오류가 발생했습니다"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "로그아웃 되었습니다"})
}

func (ac *AuthController) GetCurrentUser(c *gin.Context) {
    session := sessions.Default(c)
    userEmail := session.Get("email")
    if userEmail == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user bson.M
    err := ac.db.Collection("users").FindOne(ctx, bson.M{
        "email": userEmail.(string),
    }).Decode(&user)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류가 발생했습니다"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":    user["_id"].(primitive.ObjectID).Hex(),
        "email": user["email"],
        "name":  user["name"],
    })
}