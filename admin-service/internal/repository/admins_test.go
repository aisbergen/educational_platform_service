package repository

import (
	"admin/internal/domain"
	"admin/pkg/auth"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

type MockDB struct {
	ExecContextFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContextFunc func(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (mdb *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if mdb.ExecContextFunc != nil {
		return mdb.ExecContextFunc(ctx, query, args...)
	}

	return nil, nil
}

func (mdb *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if mdb.QueryRowContextFunc != nil {
		return mdb.QueryRowContextFunc(ctx, query, args...)
	}

	return &sql.Row{}
}

func TestCreate(t *testing.T) {
	mockDB := &MockDB{}
	repo := &AdminRepo{db: mockDB}

	admin := domain.Admin{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password",
	}

	mockDB.ExecContextFunc = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		return nil, errors.New("mock error")
	}

	err := repo.Create(context.Background(), admin)

	if err == nil {
		t.Error("expected error, got nil")
	}

	mockDB.ExecContextFunc = nil

	_ = repo.Create(context.Background(), admin)
}

func TestAdminGetByEmail(t *testing.T) {
	admin, err := createAdmin(context.Background())
	if err != nil {
		t.Error("error get by email", err)
	}
	a, err := repo.Admin.GetByEmail(context.Background(), admin.Email)
	if err != nil {
		t.Error("error get by email", err)
	}

	if a.Email != admin.Email {
		t.Error("empty email")
	}

	if a.Name != admin.Name {
		t.Error("empty name")
	}

	if a.Password != admin.Password {
		t.Error("empty password")
	}
}

func TestAdminSetSession(t *testing.T) {
	_, err := updateSetSession()
	if err != nil {
		t.Error("set session error", err)
	}
}

func TestAdminGetByRefresh(t *testing.T) {
	refresh, err := updateSetSession()
	if err != nil {
		t.Error("get by refresh token", err)
	}
	admin, err := repo.Admin.GetByRefresh(context.Background(), refresh)
	if err != nil {
		t.Error("get by refresh token", err)
	}
	if admin.Email == "" {
		t.Error("empty email")
	}
	if admin.Name == "" {
		t.Error("empty name")
	}

	if admin.Password == "" {
		t.Error("empty password")
	}

	if admin.RefreshToken == "" {
		t.Error("empty refresh")
	}

}

func createAdmin(ctx context.Context) (domain.Admin, error) {
	admin := domain.Admin{
		ID:       1,
		Email:    "test@test.com",
		Name:     "test",
		Password: "test",
	}
	err := repo.Admin.Create(ctx, admin)
	if err != nil {
		return admin, err
	}
	return admin, nil
}

func updateSetSession() (string, error) {
	manager, err := auth.NewManager(cfg.JWT.SigningKey)
	if err != nil {
		return "", err
	}
	refreshtoken, err := manager.NewRefreshToken()
	if err != nil {
		return "", err
	}
	err = repo.Admin.SetSession(context.Background(), domain.Session{
		RefreshToken: refreshtoken,
		ExpiresAt:    time.Now().Add(cfg.JWT.RefreshTokenTTL),
	}, 1)
	if err != nil {
		return "", err
	}
	return refreshtoken, nil
}
