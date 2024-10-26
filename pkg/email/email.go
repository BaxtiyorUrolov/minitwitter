package email

import (
	"net/mail"
	"net/smtp"
	"twitter/config"

	"github.com/scorredoira/email"
)

func SendEmail(to, code string) error {

	password := config.Load().EmailPassword

	from := mail.Address{Name: "BMC-Host", Address: "uralov2908@gmail.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := email.NewMessage("Verification Code", "Your verification code is: "+code)
	msg.From = from
	msg.To = []string{to}

	auth := smtp.PlainAuth("", from.Address, password, smtpHost)

	err := email.Send(smtpHost+":"+smtpPort, auth, msg)
	if err != nil {
		return err
	}
	return nil
}

func SendNotification(to, message string) error {

	password := config.Load().EmailPassword

	from := mail.Address{Name: "BMC-Host", Address: "uralov2908@gmail.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := email.NewMessage("Your tweet has been liked", "Your tweet has been liked by : "+message)
	msg.From = from
	msg.To = []string{to}

	auth := smtp.PlainAuth("", from.Address, password, smtpHost)

	err := email.Send(smtpHost+":"+smtpPort, auth, msg)
	if err != nil {
		return err
	}
	return nil
}
