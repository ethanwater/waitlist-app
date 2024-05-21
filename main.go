package main

import (
	"math/rand"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type App struct {
	Address string
	Engine  *gin.Engine
	API     *gin.RouterGroup
}

func main() {
	router := gin.Default()
	app := App{
		":8080",
		router,
		router.Group("/api"),
	}

	app.API.GET("/enrollemail", enroll())
	app.API.GET("/sendverificationcode", sendVerificationEmail())
	app.API.GET("/verifyverificationcode", verifyVerificationCode())

	app.Engine.Use(static.Serve("/", static.LocalFile("./views", true)))
	app.Engine.Run(app.Address)
}

func enroll() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

var code2FA atomic.Value

func sendVerificationEmail() func(*gin.Context) {
	return func(ctx *gin.Context) {
		adminEmail := "vivianniyuki@gmail.com"
		recipientEmail := ctx.Query("email")
		appPassword := "mmet bifn orxu uzgs"
		verificationCode, err := GenerateAuthKey2FA()
		if err != nil {
			return
		}
		hashedCode, err := HashPassword(verificationCode)
		if err != nil {
			return
		}

		m := gomail.NewMessage()
		m.SetHeader("From", adminEmail)
		m.SetHeader("To", recipientEmail)
		m.SetHeader("Subject", "Vivian's Waitlist Verification Code")
		m.SetBody("text/html", "Verification Code: "+verificationCode)

		d := gomail.NewDialer("smtp.gmail.com", 587, adminEmail, appPassword)

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		code2FA.Store(hashedCode)
		ctx.JSON(http.StatusOK, gin.H{"code": hashedCode})
	}
}

func verifyVerificationCode() func(*gin.Context) {
	return func(ctx *gin.Context) {
		inputcode := ctx.Query("code")
		hash := code2FA.Load().(string)

		result, _ := VerifyAuthKey2FA(hash, inputcode)
		if result {
			ctx.JSON(http.StatusOK, gin.H{"code": "true"})
			//store email in db
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": "false"})
		}
	}
}

const (
	charset     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	authKeySize int    = 5
)

func GenerateAuthKey2FA() (string, error) {
	source := rand.New(rand.NewSource(time.Now().Unix()))
	var authKey strings.Builder

	for i := 0; i < authKeySize; i++ {
		sample := source.Intn(len(charset))
		authKey.WriteString(string(charset[sample]))
	}
	return authKey.String(), nil
}

func VerifyAuthKey2FA(authkey_hash, input string) (bool, error) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if SanitizeCheck(input) {
		status := bcrypt.CompareHashAndPassword([]byte(authkey_hash), []byte(input))
		if status != nil {
			return status == nil, status
		} else {
			return status == nil, status
		}
	}

	return false, nil
}

const cost int = 13

func HashPassword(password string) (string, error) {
	hashChannel := make(chan struct {
		hash string
		err  error
	})
	defer close(hashChannel)

	go func() {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		hashChannel <- struct {
			hash string
			err  error
		}{string(hash), err}
	}()

	result := <-hashChannel
	return result.hash, result.err
}

func VerfiyHashPassword(hash, password string) bool {
	verificationChannel := make(chan bool)
	defer close(verificationChannel)

	go func() {
		status := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		verificationChannel <- status == nil
	}()

	return <-verificationChannel
}

var whitelist_text *regexp.Regexp = regexp.MustCompile("[^a-zA-Z0-9]+")

func Sanitize(input string) string {
	return whitelist_text.ReplaceAllString(input, "")
}

func SanitizeCheck(input string) bool {
	return whitelist_text.ReplaceAllString(input, "") == input
}

func SanitizeEmailCheck(input string) bool {
	_, err := mail.ParseAddress(input)
	return err == nil
}
