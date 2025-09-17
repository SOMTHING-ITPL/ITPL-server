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
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			fmt.Sprintf(`
				<!DOCTYPE html>
				<html lang="ko">
				<head>
					<meta charset="UTF-8">
					<style>
						body {
							font-family: Arial, Helvetica, sans-serif;
							background-color: #f9f9f9;
							padding: 20px;
						}
						.container {
							background-color: #ffffff;
							border-radius: 10px;
							padding: 20px;
							box-shadow: 0 2px 8px rgba(0,0,0,0.1);
							max-width: 400px;
							margin: auto;
						}
						.code {
							font-size: 24px;
							font-weight: bold;
							color: #2c3e50;
							background: #f0f0f0;
							padding: 10px 20px;
							border-radius: 6px;
							letter-spacing: 3px;
							text-align: center;
						}
						.footer {
							font-size: 12px;
							color: #888;
							margin-top: 20px;
							text-align: center;
						}
					</style>
				</head>
				<body>
					<div class="container">
						<h2>ITPL 이메일 인증</h2>
						<p>아래 인증 코드를 입력해주세요:</p>
						<div class="code">%s</div>
						<p class="footer">본 메일은 자동 발송되었습니다.</p>
					</div>
				</body>
				</html>
			`, code),
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
