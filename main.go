package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	app.Engine.Use(static.Serve("/", static.LocalFile("./views", true)))
	app.Engine.Run(app.Address)
}

func enroll() func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}

func sendVerificationEmail() func(*gin.Context) {
	return func(ctx *gin.Context) {
		recipientEmail := ctx.Query("email")

		_, err := mail.ParseAddress(recipientEmail)
		if err != nil {
			ctx.JSON(200, gin.H{
				"error": "invalid email",
			})
			return
		}

		adminEmail := "vivianniyuki@gmail.com"
		appPassword := "mmet bifn orxu uzgs"

		to := []string{recipientEmail}
		auth := smtp.PlainAuth("", adminEmail, appPassword, "smtp.gmail.com")

		verificationCode, err := GenerateAuthKey2FA()
		if err != nil {
			return
		}

		msg := []byte("To: " + recipientEmail + "\r\n" +
			"Subject: Vivian's Waitlist Verification Code\r\n" +
			"\r\n" +
			verificationCode +
			"\r\n")

		err = smtp.SendMail("smtp.gmail.com:587", auth, adminEmail, to, msg)
		if err != nil {
			log.Fatal(err)
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
	fmt.Println(authKey.String())

	hashChannel := make(chan string, 1)
	go func() {
		authKeyHash, err := HashPassword(authKey.String())
		if err != nil {
			hashChannel <- ""
			return
		}
		hashChannel <- authKeyHash
	}()
	hash := <-hashChannel

	if hash == "" {
		return "", nil
	}

	return hash, nil
}

func VerifyAuthKey2FA(ctx context.Context, authkey_hash, input string) (bool, error) {
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
