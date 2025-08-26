package email

import (
	"fmt"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
)

func GenerateCode(n int) string {
	digits := "0123456789"
	code := ""
	for i := 0; i < n; i++ {
		code += string(digits[rand.Intn(len(digits))])
	}
	return code
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func SendMail(targetMail string, code string) error {
	cfg := config.SmtpCfg // Smtp 설정 구조체

	msg := []byte(
		"From: " + cfg.From + "\r\n" +
			"To: " + targetMail + "\r\n" +
			"Subject: ITPL 이메일 인증 코드\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			fmt.Sprintf("인증 코드: %s\r\n", code),
	)

	// google auth
	auth := smtp.PlainAuth("", cfg.From, cfg.AppPassword, cfg.HostServer)

	// send
	err := smtp.SendMail(cfg.HostServer+":"+cfg.Port, auth, cfg.From, []string{targetMail}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent to:", targetMail)
	return nil
}
