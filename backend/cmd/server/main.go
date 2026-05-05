package main

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"skema-api/core/cache"
	"skema-api/core/config"
	"skema-api/core/db"
	"skema-api/core/mailer"
	"skema-api/core/middleware/cors"
	"skema-api/core/middleware/limiter"
	"skema-api/core/middleware/security"
	"skema-api/core/middleware/timeout"
	"skema-api/core/response"
	_ "skema-api/docs"
	"skema-api/features/auth"
)

// @title          Skema API
// @version        1.0
// @description    CMS headless self-hosted - API REST auto-générée.
// @host           localhost:3000
// @BasePath       /v1
// @securityDefinitions.apikey BearerAuth
// @in             header
// @name           Authorization

var startedAt = time.Now()

func main() {
	cfg := config.Load()
	ctx := context.Background()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	pool, err := db.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Connexion base de données échouée : %v", err)
	}
	defer pool.Close()

	c := cache.New(cfg.RedisURL)
	m := mailer.New(cfg.Mailer.From, cfg.Mailer.ResendKey)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(security.Headers(cfg.Env))
	router.Use(cors.Middleware(cfg.Security.CORSOrigins, cfg.Env))
	router.Use(limiter.Body(cfg.Security.MaxBodySizeMB))
	router.Use(limiter.Global())
	router.Use(timeout.Middleware(cfg.Security.RequestTimeoutSec))

	router.GET("/docs/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler()))
	router.GET("/favicon.ico", func(c *gin.Context) { c.Status(response.StatusNoContent) })
	router.GET("/sw.js", func(c *gin.Context) { c.Status(response.StatusNoContent) })

	api := router.Group("/v1")
	registerRoutes(api, cfg, pool, c, m)

	log.Printf("Serveur démarré sur le port %s", cfg.Port)
	log.Printf("Documentation disponible sur http://localhost:%s/docs/v1/index.html", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Impossible de démarrer le serveur : %v", err)
	}
}

func registerRoutes(api *gin.RouterGroup, cfg *config.Config, pool *pgxpool.Pool, c cache.Cache, m *mailer.Mailer) {
	api.GET("/health", func(ctx *gin.Context) { handleHealth(ctx, cfg) })

	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, c, m, cfg)
	auth.RegisterRoutes(api, authSvc, cfg.JwtSecret)
}

// @Summary      Vérification de l'état du serveur
// @Description  Retourne le statut, les infos runtime et les fonctionnalités disponibles.
// @Tags         système
// @Produce      json
// @Success      200  {object}  response.Body
// @Router       /health [get]

func handleHealth(c *gin.Context, cfg *config.Config) {
	response.OK(c, "Serveur opérationnel.", gin.H{
		"version":         APIVersion,
		"env":             cfg.Env,
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
		"uptime":          time.Since(startedAt).Round(time.Second).String(),
		"fonctionnalités": availableFeatures,
		"runtime": gin.H{
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
		},
	})
}
