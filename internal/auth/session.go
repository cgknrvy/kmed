package auth

import (
	"context"
	"time"

	"kmed/api/ent"
	entSession "kmed/api/ent/session"

	"github.com/google/uuid"
)

type SessionRepo struct {
	client *ent.Client
}

func (s SessionRepo) Create(
	ctx context.Context,
	userID uuid.UUID,
	hash string,
	expiresAt time.Time,
) error {
	_, err := s.client.Session.Create().
		SetUserID(userID).
		SetRefreshTokenHash(hash).
		SetExpiresAt(expiresAt).
		Save(ctx)
	return err
}

func (s SessionRepo) FindByHash(ctx context.Context, hash string) (*ent.Session, error) {
	session, err := s.client.Session.Query().Where(entSession.RefreshTokenHashEQ(hash)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	return s.client.Session.Update().
		Where(entSession.And(entSession.IDEQ(id), entSession.RevokedAtIsNil())).
		SetRevokedAt(time.Now()).
		Exec(ctx)
}

func (s *SessionRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return s.client.Session.Update().
		Where(entSession.And(entSession.UserID(userID), entSession.RevokedAtIsNil())).
		SetRevokedAt(time.Now()).
		Exec(ctx)
}

func (s *SessionRepo) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := s.client.Session.Delete().
		Where(entSession.And(entSession.UserID(userID), entSession.RevokedAtIsNil())).
		Exec(ctx)
	return err
}
