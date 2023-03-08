package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func Logger(l *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			l.Info(r.URL.String())

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
