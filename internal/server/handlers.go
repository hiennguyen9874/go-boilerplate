package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/hiennguyen9874/go-boilerplate/config"
	itemHttp "github.com/hiennguyen9874/go-boilerplate/internal/items/delivery/http"
	itemRepository "github.com/hiennguyen9874/go-boilerplate/internal/items/repository"
	itemUseCase "github.com/hiennguyen9874/go-boilerplate/internal/items/usecase"
	apiMiddleware "github.com/hiennguyen9874/go-boilerplate/internal/middleware"
	userHttp "github.com/hiennguyen9874/go-boilerplate/internal/users/delivery/http"
	userRepository "github.com/hiennguyen9874/go-boilerplate/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/go-boilerplate/internal/users/usecase"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Repository
	userPgRepo := userRepository.CreateUserPgRepository(db)
	userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)
	itemPgRepo := itemRepository.CreateItemPgRepository(db)

	// UseCase
	userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, cfg, logger)
	itemUC := itemUseCase.CreateItemUseCaseI(itemPgRepo, cfg, logger)

	// Handler
	userHandler := userHttp.CreateUserHandler(userUC, cfg, logger)
	itemHandler := itemHttp.CreateItemHandler(itemUC, cfg, logger)

	// middleware
	mw := apiMiddleware.CreateMiddlewareManager(cfg, logger, userUC)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(time.Duration(cfg.Server.ProcessTimeout) * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(mw.Cors()))

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	userHttp.MapUserRoute(apiRouter, userHandler, mw)
	itemHttp.MapItemRoute(apiRouter, itemHandler, mw)

	return r, nil
}
