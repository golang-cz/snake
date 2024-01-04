package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/golang-cz/snake/proto"
)

const (
	PORT                 = 5252
	NumOfAISnakes        = 3
	FoodFromDeadSnake    = 3
	GameTickTime         = 75 * time.Millisecond
	NumOfStartingFood    = 3
	FoodGenerateInterval = 2 * time.Second
	AISnakeRespawnTime   = 5 * time.Second
)

func main() {
	slog.Info(fmt.Sprintf("serving at http://localhost:%v", PORT))

	rpc := NewSnakeServer()
	go rpc.Run(context.TODO())

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", PORT), rpc.Router()); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	// r.Use(requestDebugger)
	// r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	webrpcHandler := proto.NewSnakeGameServer(s)
	r.Handle("/*", webrpcHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	return r
}

func requestDebugger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("req started",
			slog.String("url", fmt.Sprintf("%v %v", r.Method, r.URL.String())))

		defer func() {
			slog.Info("req finished",
				slog.String("url", fmt.Sprintf("%v %v", r.Method, r.URL.String())),
			)
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
