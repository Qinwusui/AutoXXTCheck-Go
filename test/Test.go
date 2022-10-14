package main

import "gopkg.in/gomail.v2"

func main() {
	sendEmail("2064508450","？？？")	
}

func sendEmail(mailTo string, msg string) (error, bool) {
	mailConn := map[string]string{
		"user":     "loverqinthisway@outlook.com",
		"password": "Xyz5201314",
		"host":     "smtp.office365.com",
	}
	port := 587
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "五岁"))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", "学习通打卡回执")
	m.SetBody("text/html", `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h2>Maser</h2>
    <p>你的注册验证码为`+msg+`</p>
	<p>请妥善保管</p>
</body>
</html>
`)
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["password"])
	err := d.DialAndSend(m)
	return err, true
}