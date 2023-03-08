package server

import (
	"context"
	"github.com/Silarsone/image-loader/internal/config"
	"github.com/Silarsone/image-loader/internal/mock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestServer_HandleGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockController := mock.NewMockcontroller(ctrl)

	mockController.EXPECT().GetUser(gomock.AssignableToTypeOf(context.Background()), int64(1))

	keyword := "mykeyword"

	s := NewServer("", logrus.New(), mockController, &config.Config{JWTKeyword: keyword})
	s.RegisterRoutes()

	srv := httptest.NewServer(s.r)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/user/1", nil)
	require.NoError(t, err)

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(1),
		Subject:   "authorized",
		Audience:  []string{"1"},
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        strconv.Itoa(1),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(keyword))
	require.NoError(t, err)

	req.Header.Set("Authorization", tokenString)

	client := http.Client{
		Timeout: time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
