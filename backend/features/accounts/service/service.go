package service

import (
	"skema-api/core/cache"
	"skema-api/core/config"
	"skema-api/core/mailer"
	"skema-api/features/accounts/repository"
)

type Service struct {
	repo        *repository.Repository
	cache       cache.Cache
	mailer      *mailer.Mailer
	jwtSecret   string
	frontendURL string
}

/*
 * New instancie le service d'authentification avec ses dépendances.
 *
 * Attend  : le repository, le cache, le mailer et la configuration globale.
 * Retourne: un pointeur vers Service prêt à l'emploi.
 */

func New(repo *repository.Repository, c cache.Cache, m *mailer.Mailer, cfg *config.Config) *Service {
	return &Service{
		repo:        repo,
		cache:       c,
		mailer:      m,
		jwtSecret:   cfg.JwtSecret,
		frontendURL: cfg.FrontendURL,
	}
}
