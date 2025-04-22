package utils

import (
	"fmt"
	"net/smtp"
)

func SendVerificationEmail(email string, token string) (bool, error) {
	// Генерация 32 случайных байт (256 бит)
	from := "allusion@debugmail.io"
	host := "app.debugmail.io"
	addr := "app.debugmail.io:25"

	login := "d1ebb98a-447b-4488-bc8e-ad3962683fb0"
	pass := "6f547321-f68f-497a-a771-d0e162c4e503"
	subject := "Email confirmation"
	auth := smtp.PlainAuth("", login, pass, host)
	to := []string{email}
	// msg := []byte("Subject: Hello! This is the body of the email.")

	verificationLink := "http://localhost:8080/account/verify/" + token

	// Тело письма в формате MIME
	body := fmt.Sprintf(`To: %s
From: %s
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset="UTF-8"

<html>
<body>
    <h1>Подтвердите ваш email</h1>
    <p>Нажмите <a href="%s">здесь</a> для завершения регистрации</p>
    <p>Или скопируйте ссылку: %s</p>
</body>
</html>`, email, from, subject, verificationLink, verificationLink)

	msg := []byte(body)

	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	fmt.Println("Email sent successfully!")
	return true, nil
}
