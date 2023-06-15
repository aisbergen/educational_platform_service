package service

import (
	"admin/internal/domain"
	repoMock "admin/internal/repository/mocks"
	"admin/pkg/auth"
	"admin/pkg/hash"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAdminService_SignUp(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockAdmin(ctl)

	ctx := context.Background()
	admin := domain.Admin{
		Email:    "test1@test.com",
		Name:     "test",
		Password: "test2001",
	}

	hash := hash.NewHash(10)

	repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	manager, _ := auth.NewManager("qwerty")

	service := NewAdminService(repo, hash, manager, 15*time.Minute, 15*time.Minute)

	err := service.SignUp(ctx, admin)
	require.NoError(t, err)
}

func TestAdminService_SignUpError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockAdmin(ctl)

	ctx := context.Background()

	hash := hash.NewHash(10)

	repo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

	manager, _ := auth.NewManager("qwerty")

	service := NewAdminService(repo, hash, manager, 15*time.Minute, 15*time.Minute)

	err := service.SignUp(ctx, domain.Admin{})
	require.Nil(t, err)
}

func TestAdminService_SignIn(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockAdmin(ctl)

	ctx := context.Background()
	email := "test1@test.com"
	password := "test2001"

	hash := hash.NewHash(10)
	password, err := hash.GenerateFromPassword("test2001")
	require.NoError(t, err)
	admin := domain.Admin{
		ID:       1,
		Email:    "test1@test.com",
		Name:     "test",
		Password: password,
	}

	manager, _ := auth.NewManager("qwerty")

	repo.EXPECT().GetByEmail(ctx, gomock.Any()).Return(admin, nil).Times(1)
	repo.EXPECT().SetSession(ctx, gomock.Any(), admin.ID).Return(nil).Times(1)

	service := NewAdminService(repo, hash, manager, 15*time.Minute, 15*time.Minute)
	_, err = service.SignIn(ctx, email, "test2001")
	require.NoError(t, err)
}
