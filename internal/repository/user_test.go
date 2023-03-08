package repository

import (
	"context"
	"fmt"
	"github.com/Silarsone/image-loader/internal/config"
	"github.com/Silarsone/image-loader/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserRepo_AddUser(t *testing.T) {
	cfg := config.Config{DB: &config.DB{
		Driver:   "postgres",
		Password: "secretpassword",
		User:     "postgres",
		Name:     "postgres",
		SSLMode:  "disable",
	}}

	db, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s", cfg.DB.User,
		cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Password))
	require.NoError(t, err)

	userRepo := NewUserRepo(db, cfg.DB)

	err = userRepo.AddUser(context.Background(), model.User{
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	})
	assert.NoError(t, err)

	user, err := userRepo.CheckAuth(context.Background(), "test", "test")
	assert.Equal(t, "test", user.Login)
	assert.Equal(t, "test", user.Password)
}
