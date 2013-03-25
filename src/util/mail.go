package util

import (
	"net/smtp"
	"strings"
)

/*
usage:
	user := "xxxx@163.com"
	password := "xxxx"
	host := "smtp.163.com:25"
	to := "xxxx@gmail.com;ssssss@gmail.com"

	subject := "Test send email by golang"

	body := `
	<html>
	<body>
	<h3>
	"Test send email by golang"
	</h3>
	</body>
	</html>
	`
	fmt.Println("send email")
	err := SendMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("send mail error!")
		fmt.Println(err)
	}else{
		fmt.Println("send mail success!")
	}

*/
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := &unencryptedAuth{
		smtp.PlainAuth(
			"",
			user,
			password,
			hp[0],
		),
	}
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

/*
为了解决"unencrypted connection"问题
这个问题是因为TLS没有返回true

*/
type unencryptedAuth struct {
	smtp.Auth
}

func (a *unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
