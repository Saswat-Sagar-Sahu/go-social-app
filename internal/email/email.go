package email

import "context"

// Sender is the interface used by the application to send emails.
type Sender interface {
    Send(ctx context.Context, to, subject, plainText, htmlBody string) error
}
