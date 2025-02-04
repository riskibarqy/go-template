package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/riskibarqy/99backend-challenge/config"
	"github.com/riskibarqy/99backend-challenge/internal/data"
	"github.com/riskibarqy/99backend-challenge/internal/http/controller"
	"github.com/riskibarqy/99backend-challenge/internal/user"
	"github.com/rs/cors"
)

// Server represents the http server that handles the requests
type Server struct {
	dataManager    *data.Manager
	userService    user.ServiceInterface
	userController *controller.UserController
}

func (hs *Server) authMethod(r chi.Router, method string, path string, handler http.HandlerFunc) {
	r.With(
		hs.instrument(method, "/v1"+path),
	).Method(method, path, handler)
}

func (hs *Server) compileRouter() chi.Router {
	r := chi.NewRouter()

	// Base middlewares
	//

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	//Routes()
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Access-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// Add routes

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	r.HandleFunc("/public-api/login", hs.userController.Login)

	// r.Route("/users", func(r chi.Router) {
	// 	r.Post("/", hs.userController.CreateUser)
	// 	r.Get("/", hs.userController.ListUser)
	// 	r.Get("/{userId}", hs.userController.GetUserByID)
	// })

	r.Route("/public-api/users", func(r chi.Router) {
		r.Post("/", hs.userController.CreateUser)
	})

	r.Route("/private", func(r chi.Router) {
		r.Use(hs.authorizedOnly(hs.userService))

		// private users
		hs.authMethod(r, "GET", "/logout", hs.userController.Logout)
		hs.authMethod(r, "PUT", "/users/changePassword", hs.userController.ChangePassword)
		hs.authMethod(r, "PUT", "/users/{userId}", hs.userController.UpdateUser)
		hs.authMethod(r, "GET", "/users", hs.userController.ListUser)
		hs.authMethod(r, "GET", "/users/{userId}", hs.userController.GetUserByID)
		hs.authMethod(r, "POST", "/users", hs.userController.CreateUser)
	})

	return r
}

// Serve serves http requests
func (hs *Server) Serve() {
	// Compile all the routes
	//

	r := hs.compileRouter()

	// Run the server + gracefully shutdown mechanism
	//

	log.Printf("About to listen on 8080. Go to http://127.0.0.1:8080")
	srv := http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// NewServer creates a new http server
func NewServer(
	config *config.Config,
	dataManager *data.Manager,
	userService user.ServiceInterface,
) *Server {
	userController := controller.NewUserController(userService, dataManager)

	return &Server{
		dataManager:    dataManager,
		userService:    userService,
		userController: userController,
	}
}
