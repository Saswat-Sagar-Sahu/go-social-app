package email

import (
    "context"
    "fmt"

    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"

    "github.com/Saswat-Sagar-Sahu/Social/internal/env"
)

type SendGridSender struct {
    client    *sendgrid.Client
    fromEmail string
    fromName  string
}

func NewSendGridSender(apiKey, fromEmail, fromName string) *SendGridSender {
    return &SendGridSender{
        client:    sendgrid.NewSendClient(apiKey),
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

// NewSendGridFromEnv constructs a SendGrid sender using environment variables.
// Expects SENDGRID_API_KEY and INVITE_FROM (email). INVITE_FROM_NAME is optional.
func NewSendGridFromEnv() (*SendGridSender, error) {
    apiKey := env.GetString("SENDGRID_API_KEY", "")
    if apiKey == "" {
        return nil, fmt.Errorf("SENDGRID_API_KEY not set")
    }
    from := env.GetString("INVITE_FROM", "no-reply@example.com")
    fromName := env.GetString("INVITE_FROM_NAME", "Social App")
    return NewSendGridSender(apiKey, from, fromName), nil
}

func (s *SendGridSender) Send(ctx context.Context, to, subject, plainText, htmlBody string) error {
    from := mail.NewEmail(s.fromName, s.fromEmail)
    toEmail := mail.NewEmail("", to)
    message := mail.NewSingleEmail(from, subject, toEmail, plainText, htmlBody)

    resp, err := s.client.Send(message)
    if err != nil {
        return err
    }
    if resp.StatusCode >= 400 {
        return fmt.Errorf("sendgrid send failed: status=%d body=%s", resp.StatusCode, resp.Body)
    }
    return nil
}
