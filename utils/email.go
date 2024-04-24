package utils

import (
	"github.com/spf13/viper"
	"net/smtp"
)

func SendEMail(toEmail string, boBy string) error {
	auth := smtp.PlainAuth(
		"",
		viper.GetString("smtp.username"),
		viper.GetString("smtp.password"),
		viper.GetString("smtp.host"),
	)
	sendAutor := viper.GetString("smtp.username")
	msg := []byte("From:" + sendAutor + "\r\n" + "To: " + toEmail + "\r\n" + "Subject: 告警信息\r\n" + "\r\n" + "内容：" + boBy + ".\r\n")
	err := smtp.SendMail(
		viper.GetString("smtp.addr"),
		// "smtp.qq.com:587",
		auth,
		toEmail,
		[]string{toEmail},
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}
