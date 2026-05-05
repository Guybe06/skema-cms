package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
	DelPattern(ctx context.Context, pattern string) error
}

/*
 * New instancie le cache Redis si l'URL est fournie, sinon bascule sur le cache mémoire.
 *
 * Attend  : l'URL Redis (peut être vide pour le mode mémoire).
 * Retourne: une implémentation de Cache prête à l'emploi.
 */

func New(redisURL string) Cache {
	if redisURL == "" {
		return newMemoryCache()
	}

	client, err := newRedisCache(redisURL)
	if err != nil {
		return newMemoryCache()
	}

	return client
}
