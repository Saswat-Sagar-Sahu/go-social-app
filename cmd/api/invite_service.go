package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/Saswat-Sagar-Sahu/Social/internal/env"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
)

// CreateAndSendInvite creates an invitation record and sends the email if configured.
func (app *application) CreateAndSendInvite(ctx context.Context, inviterID *int64, emailAddr string, message *string) (*store.Invitation, error) {
	// generate secure token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(b)

	expiryHours := env.GetInt("INVITE_TOKEN_EXPIRY_HOURS", 48)

	inv := &store.Invitation{
		InviterID:    inviterID,
		InviteeEmail: emailAddr,
		Message:      message,
		Token:        token,
		ExpiresAt:    time.Now().Add(time.Duration(expiryHours) * time.Hour),
	}

	if err := app.Store.Invitations.Create(ctx, inv); err != nil {
		return nil, err
	}

	if app.Email != nil {
		acceptURL := fmt.Sprintf("%s/invites/accept?token=%s", app.config.apiURL, token)
		subject := "You were invited to join"
		plain := fmt.Sprintf("You were invited to join.\n\nMessage: %s\n\nAccept: %s", func() string {
			if message != nil {
				return *message
			}
			return ""
		}(), acceptURL)
		html := fmt.Sprintf("<p>You were invited to join.</p><p>%s</p><p><a href=\"%s\">Accept invite</a></p>", func() string {
			if message != nil {
				return *message
			}
			return ""
		}(), acceptURL)

		if err := app.Email.Send(ctx, emailAddr, subject, plain, html); err != nil {
			log.Println("failed to send invite email:", err)
		}
	}

	return inv, nil
}
