package store

import (
	"context"
	"database/sql"
	"time"
)

type Invitation struct {
	ID           string    `json:"id"`
	InviterID    *int64    `json:"inviter_id,omitempty"`
	InviteeEmail string    `json:"invitee_email"`
	Message      *string   `json:"message,omitempty"`
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    string    `json:"created_at"`
	AcceptedAt   *string   `json:"accepted_at,omitempty"`
	RedeemedBy   *int64    `json:"redeemed_by,omitempty"`
}

type InvitationsStore struct {
	db *sql.DB
}

func (s *InvitationsStore) Create(ctx context.Context, inv *Invitation) error {
	query := `INSERT INTO invitations (inviter_id, invitee_email, message, token, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	var inviterID interface{}
	if inv.InviterID != nil {
		inviterID = *inv.InviterID
	} else {
		inviterID = nil
	}
	err := s.db.QueryRowContext(ctx, query, inviterID, inv.InviteeEmail, inv.Message, inv.Token, inv.ExpiresAt).Scan(&inv.ID, &inv.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *InvitationsStore) GetByToken(ctx context.Context, token string) (*Invitation, error) {
	query := `SELECT id, inviter_id, invitee_email, message, token, expires_at, created_at, accepted_at, redeemed_by FROM invitations WHERE token = $1`
	inv := &Invitation{}
	var inviter sql.NullInt64
	var message sql.NullString
	var accepted sql.NullString
	var redeemed sql.NullInt64
	err := s.db.QueryRowContext(ctx, query, token).Scan(
		&inv.ID,
		&inviter,
		&inv.InviteeEmail,
		&message,
		&inv.Token,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&accepted,
		&redeemed,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	if inviter.Valid {
		v := inviter.Int64
		inv.InviterID = &v
	}
	if message.Valid {
		m := message.String
		inv.Message = &m
	}
	if accepted.Valid {
		a := accepted.String
		inv.AcceptedAt = &a
	}
	if redeemed.Valid {
		r := redeemed.Int64
		inv.RedeemedBy = &r
	}
	return inv, nil
}

func (s *InvitationsStore) MarkAccepted(ctx context.Context, token string, redeemedBy int64) error {
	query := `UPDATE invitations SET accepted_at = $1, redeemed_by = $2 WHERE token = $3`
	_, err := s.db.ExecContext(ctx, query, time.Now(), redeemedBy, token)
	return err
}
