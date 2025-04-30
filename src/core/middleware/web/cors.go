package web

import (
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golibs-starter/golib"
	"go.uber.org/fx"
	"net/http"
)

func CORSOpt() fx.Option {
	return fx.Options(
		fx.Invoke(RegisterCORS),
	)
}

type RegisterCORSIn struct {
	fx.In
	App       *golib.App
	CORSProps *properties.CORSProperties
}

func RegisterCORS(in RegisterCORSIn) {
	in.App.AddHandler(
		CORSHandler,
	)
}
func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If preflight, respond directly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call next handler
		next.ServeHTTP(w, r)
	})
}
