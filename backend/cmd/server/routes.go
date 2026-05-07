package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"runtime"
	"time"

	"skema-api/core/cache"
	"skema-api/core/config"
	"skema-api/core/mailer"
	"skema-api/core/response"
	"skema-api/features/accounts"
	accountsrepo "skema-api/features/accounts/repository"
	accountssvc "skema-api/features/accounts/service"
	"skema-api/features/apikeys"
	apikeysrepo "skema-api/features/apikeys/repository"
	apikeyssvc "skema-api/features/apikeys/service"
	"skema-api/features/collections"
	collectionsrepo "skema-api/features/collections/repository"
	collectionssvc "skema-api/features/collections/service"
	"skema-api/features/connections"
	connectionsrepo "skema-api/features/connections/repository"
	connectionssvc "skema-api/features/connections/service"
	"skema-api/features/content"
	contentsvc "skema-api/features/content/service"
	"skema-api/features/memberships"
	membershipsrepo "skema-api/features/memberships/repository"
	membershipssvc "skema-api/features/memberships/service"
	"skema-api/features/organizations"
	orgsrepo "skema-api/features/organizations/repository"
	orgssvc "skema-api/features/organizations/service"
	"skema-api/features/publicapi"
	"skema-api/features/users"
	usersrepo "skema-api/features/users/repository"
	userssvc "skema-api/features/users/service"
)

func registerRoutes(api *gin.RouterGroup, cfg *config.Config, pool *pgxpool.Pool, c cache.Cache, m *mailer.Mailer) {
	api.GET("/health", func(ctx *gin.Context) { handleHealth(ctx, cfg) })

	repo := accountsrepo.New(pool)
	svc := accountssvc.New(repo, c, m, cfg)
	accounts.RegisterRoutes(api, svc, cfg.JwtSecret)

	users.RegisterRoutes(api, userssvc.New(usersrepo.New(pool)), cfg.JwtSecret)

	orgsRepository := orgsrepo.New(pool)
	organizations.RegisterRoutes(api, orgssvc.New(orgsRepository), cfg.JwtSecret)
	memberships.RegisterRoutes(api, membershipssvc.New(membershipsrepo.New(pool), orgsRepository, m, cfg.FrontendURL), cfg.JwtSecret)

	connSvc := connectionssvc.New(connectionsrepo.New(pool), orgsRepository, cfg.EncryptionKey)
	connections.RegisterRoutes(api, connSvc, cfg.JwtSecret)

	collRepo := collectionsrepo.New(pool)
	collections.RegisterRoutes(api, collectionssvc.New(collRepo, orgsRepository, connSvc), cfg.JwtSecret)
	content.RegisterRoutes(api, contentsvc.New(collRepo, orgsRepository, connSvc), cfg.JwtSecret)

	apikeySvc := apikeyssvc.New(apikeysrepo.New(pool), orgsRepository)
	apikeys.RegisterRoutes(api, apikeySvc, cfg.JwtSecret)

	pub := api.Group("/pub")
	publicapi.RegisterRoutes(pub, apikeySvc, orgsRepository, collRepo, connSvc)
}

// @Summary      Vérification de l'état du serveur
// @Tags         système
// @Produce      json
// @Success      200  {object}  response.Body
// @Router       /health [get]
func handleHealth(c *gin.Context, cfg *config.Config) {
	response.OK(c, MsgHealthOK, gin.H{
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
