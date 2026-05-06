package service

import (
	"skema-api/core/mailer"
	"skema-api/features/memberships/repository"
	orgsrepo "skema-api/features/organizations/repository"
)

type Service struct {
	repo        *repository.Repository
	orgsRepo    *orgsrepo.Repository
	mailer      *mailer.Mailer
	frontendURL string
}

func New(repo *repository.Repository, orgsRepo *orgsrepo.Repository, m *mailer.Mailer, frontendURL string) *Service {
	return &Service{repo: repo, orgsRepo: orgsRepo, mailer: m, frontendURL: frontendURL}
}
