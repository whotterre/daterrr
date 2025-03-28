package utils

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name       string
	sourceAddr string
	sourcePwd  string
}

// NewGmailSender creates a sender for emails using Gmail
func NewGmailSender(name, sourceAddr, sourcePwd string) EmailSender {
	return &GmailSender{
		name:       name,
		sourceAddr: sourceAddr,
		sourcePwd:  sourcePwd,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.sourceAddr)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	// Attach files if any
	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}
	
	smtpAuth := smtp.PlainAuth("", sender.sourceAddr, sender.sourcePwd, smtpAuthAddress)

	if err := e.Send(smtpServerAddress, smtpAuth); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

