package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-cz/snake/proto"
)

func main() {
	port := 5252
	slog.Info(fmt.Sprintf("serving at http://localhost:%v", port))

	rpc := NewSnakeServer()
	go rpc.Run(context.TODO())

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), rpc.Router())
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(requestDebugger)
	//r.Use(middleware.Recoverer)

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
		slog.Info(fmt.Sprintf("req started"),
			slog.String("url", fmt.Sprintf("%v %v", r.Method, r.URL.String())))

		defer func() {
			slog.Info(fmt.Sprintf("req finished"),
				slog.String("url", fmt.Sprintf("%v %v", r.Method, r.URL.String())),
			)
		}()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
