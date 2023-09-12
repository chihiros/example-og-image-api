package infra

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chihiros/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	_ "embed"
)

//go:embed embed/NotoSansJP-Medium.otf
var fontTitle []byte

//go:embed embed/NotoSansJP-Regular.otf
var fontUserName []byte

//go:embed embed/logo.png
var logo []byte

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.Logger)
	r.Use(middleware.Recoverer)

	// Access-Control-Allow-Originを許可する
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// OG画像用のAPI
	c := NewController()
	r.Route("/og", func(r chi.Router) {
		r.Get("/", c.GenOgImage)
	})

	// 疎通確認用のAPI
	r.Route("/now", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			jst, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				panic(err)
			}
			now := time.Now().In(jst)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(now)
		})
	})

	return r
}
