package mail

import (
	"log"
	"net/smtp"
)

func sendMail() {
	// 设置认证信息。
	auth := smtp.PlainAuth(
		"",
		"user@example.com",
		"password",
		"mail.example.com",
	)
	// 连接到服务器, 认证, 设置发件人、收件人、发送的内容,
	// 然后发送邮件。
	err := smtp.SendMail(
		"mail.example.com:25",
		auth,
		"sender@example.org",
		[]string{"recipient@example.net"},
		[]byte("To: recipient@example.net\r\n"+
			"From: sender@example.org\r\n"+
			"Subject: 邮件主题\r\n"+
			"Content-Type: text/plain; "+
			"charset=UTF-8\r\n\r\nHello World"),
	)
	if err != nil {
		log.Fatal(err)
	}

}
