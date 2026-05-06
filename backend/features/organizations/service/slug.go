package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func (s *Service) buildUniqueSlug(ctx context.Context, name string) (string, error) {
	base := nonAlphanumeric.ReplaceAllString(strings.ToLower(name), "-")
	base = strings.Trim(base, "-")

	slug := base
	for i := 2; ; i++ {
		exists, err := s.repo.SlugExists(ctx, slug)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}
