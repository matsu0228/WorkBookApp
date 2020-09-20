package api

import (
	"net/smtp"
)

type mail struct {
	from     string
	username string
	password string
	to       string
	sub      string
	msg      string
}

func (m mail) body() string {
	return "To: " + m.to + "\r\n" +
		"Subject: " + m.sub + "\r\n\r\n" +
		m.msg + "\r\n"

}

func gmailSend(send string) error {
	m := mail{
		from:     "workbook.app.golang@gmail.com",
		username: "workbook.app.golang@gmail.com",
		password: "Goland6028",
		to:       send,
		msg:      "パスワード再発行は下記のurlからお願い致します。" + "\r\n" + "http://localhost:8080/login/recover-password",
	}

	smtpSvr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", m.username, m.password, "smtp.gmail.com")
	if err := smtp.SendMail(smtpSvr, auth, m.from, []string{m.to}, []byte(m.body())); err != nil {
		return err
	}
	return nil
}
