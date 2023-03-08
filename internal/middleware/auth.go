package middleware

import (
	"context"
	"github.com/Silarsone/image-loader/internal/constants"
	"github.com/Silarsone/image-loader/internal/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func Auth(keyword string, l *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(keyword), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)

				b, err := response.ParseResponse("unauthorized", true)
				if err != nil {
					l.Error(err)
					return
				}

				_, err = w.Write(b)
				if err != nil {
					l.Error(err)
					return
				}
				return
			}

			idStr, err := token.Claims.GetIssuer()
			if err != nil {
				writeErr(err, l, w)
				return
			}

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				writeErr(err, l, w)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.IdCtxKey, int(id))

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func writeErr(err error, l *logrus.Logger, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	l.Error(err)

	b, err := response.ParseResponse(err.Error(), true)
	if err != nil {
		l.Error(err)
	}

	_, err = w.Write(b)
	if err != nil {
		l.Error(err)
	}
}
