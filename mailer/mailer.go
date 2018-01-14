package mailer

import (
	"fmt"
	"net/smtp"
	"os"

	"../parser"
)

// Send receives a recipient and body and delivers the email
func Send(jobs []parser.Job) error {
	from := os.Getenv("SEEKER_MAILER_USER")
	to := os.Getenv("SEEKER_MAILER_TO")
	auth := smtp.PlainAuth("", from, os.Getenv("SEEKER_MAILER_PASS"), "smtp.gmail.com")
	body := buildBody(jobs)
	msg := "From: <Seeker>" + from + "\n" +
		"To: " + to + "\n" +
		"Subject: New jobs from Seeker\n\n" + body
	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		to,
		[]string{to},
		[]byte(msg),
	)
}

func buildBody(jobs []parser.Job) string {
	var body string
	for _, job := range jobs {
		body += fmt.Sprintf("%s (%s) %s\n\n", job.Company, job.Title, job.Link)
	}
	return body
}
