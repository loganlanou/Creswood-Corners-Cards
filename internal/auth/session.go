package auth

import (
	"context"
	"net/http"
	"time"

	"creswoodcornerscards/internal/config"
	"creswoodcornerscards/internal/data"
)

type Manager struct {
	cfg   config.Config
	store *data.Store
}

func NewManager(cfg config.Config, store *data.Store) *Manager {
	return &Manager{cfg: cfg, store: store}
}

func (m *Manager) CurrentUser(r *http.Request) (*data.User, error) {
	cookie, err := r.Cookie(m.cfg.SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, nil
	}
	user, err := m.store.FindUserBySession(r.Context(), cookie.Value)
	if err != nil {
		return nil, nil
	}
	return &user, nil
}

func (m *Manager) SignIn(ctx context.Context, w http.ResponseWriter, userID int64) error {
	token, err := m.store.CreateSession(ctx, userID)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     m.cfg.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
	return nil
}

func (m *Manager) SignOut(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(m.cfg.SessionCookieName)
	if err == nil && cookie.Value != "" {
		_ = m.store.DeleteSession(ctx, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     m.cfg.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	return nil
}
