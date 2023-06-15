package repository

import (
	"admin/internal/domain"
	"context"
	"database/sql"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Admin interface {
	Create(ctx context.Context, admin domain.Admin) error
	GetByEmail(ctx context.Context, email string) (domain.Admin, error)
	SetSession(ctx context.Context, session domain.Session, id int) error
	GetByRefresh(ctx context.Context, refreshToken string) (domain.Admin, error)
}

type Repository struct {
	Admin Admin
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Admin: NewAdminRepo(db),
	}
}
