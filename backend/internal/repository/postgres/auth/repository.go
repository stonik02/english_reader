package auth

import (
	"context"
	"errors"
	"fmt"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool} }

func (r *Repository) CreateUserWithSession(c context.Context, email, pass string, h []byte, exp time.Time, device string) (domain.User, error) {
	tx, e := r.pool.Begin(c)
	if e != nil {
		return domain.User{}, e
	}

	defer tx.Rollback(c)

	u := domain.User{ID: uuid.NewString(), Email: email, CreatedAt: time.Now()}

	_, e = tx.Exec(c, `INSERT INTO users (id,email,password_hash) VALUES ($1,$2,$3)`, u.ID, email, pass)
	if e != nil {
		var p *pgconn.PgError
		if errors.As(e, &p) && p.Code == "23505" {
			return domain.User{}, domain.ErrEmailTaken
		}
		return domain.User{}, fmt.Errorf("insert user: %w", e)
	}

	_, e = tx.Exec(c, `INSERT INTO auth_sessions (id,user_id,refresh_token_hash,device_label,expires_at) VALUES ($1,$2,$3,$4,$5)`, uuid.NewString(), u.ID, h, device, exp)
	if e != nil {
		return domain.User{}, e
	}

	e = tx.Commit(c)

	return u, e
}
func (r *Repository) FindUserByEmail(c context.Context, email string) (domain.User, string, error) {
	var u domain.User
	var p string

	e := r.pool.QueryRow(c, `SELECT id,email,created_at,password_hash FROM users WHERE email=$1 AND deleted_at IS NULL`, email).
		Scan(&u.ID, &u.Email, &u.CreatedAt, &p)
	if errors.Is(e, pgx.ErrNoRows) {
		return domain.User{}, "", domain.ErrInvalidCredentials
	}

	return u, p, e
}
func (r *Repository) FindUserByID(c context.Context, id string) (domain.User, error) {
	var u domain.User

	e := r.pool.QueryRow(c, `SELECT id,email,created_at FROM users WHERE id=$1 AND deleted_at IS NULL`, id).
		Scan(&u.ID, &u.Email, &u.CreatedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrInvalidCredentials
	}

	return u, e
}
func (r *Repository) FindUserByRefreshHash(c context.Context, h []byte) (domain.User, error) {
	var u domain.User

	e := r.pool.QueryRow(c, `
	SELECT u.id,u.email,u.created_at FROM auth_sessions s 
    JOIN users u ON u.id=s.user_id 
	    WHERE s.refresh_token_hash=$1 AND s.revoked_at IS NULL 
	      AND s.expires_at>NOW() AND u.deleted_at IS NULL
	      `, h).
		Scan(&u.ID, &u.Email, &u.CreatedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrInvalidCredentials
	}

	return u, e
}
func (r *Repository) CreateSession(c context.Context, id string, h []byte, exp time.Time, d string) error {
	_, e := r.pool.Exec(c,
		`INSERT INTO auth_sessions (id,user_id,refresh_token_hash,device_label,expires_at) VALUES ($1,$2,$3,$4,$5)`,
		uuid.NewString(), id, h, d, exp)

	return e
}
func (r *Repository) RevokeSession(c context.Context, h []byte) error {
	_, e := r.pool.Exec(c, `UPDATE auth_sessions SET revoked_at=NOW() WHERE refresh_token_hash=$1 AND revoked_at IS NULL`, h)
	return e
}
