package conduit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type poolEntry struct {
	conn     Conduit
	lastUsed time.Time
}

/*
 * Pool gère un ensemble de connexions Conduit par identifiant (connectionID).
 *
 * Attend  : un Factory capable d'instancier un Conduit depuis un driver et une DSN.
 * Retourne: la connexion existante ou en crée une nouvelle si absente ou expirée.
 */

type Pool struct {
	mu      sync.RWMutex
	entries map[string]*poolEntry
	factory Factory
}

type Factory func(ctx context.Context, driver, dsn string) (Conduit, error)

func NewPool(factory Factory) *Pool {
	p := &Pool{
		entries: make(map[string]*poolEntry),
		factory: factory,
	}
	go p.cleanup()
	return p
}

func (p *Pool) Get(ctx context.Context, id, driver, dsn string) (Conduit, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if entry, ok := p.entries[id]; ok {
		if err := entry.conn.Ping(ctx); err == nil {
			entry.lastUsed = time.Now()
			return entry.conn, nil
		}
		_ = entry.conn.Close()
		delete(p.entries, id)
	}

	conn, err := p.factory(ctx, driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", ErrConnectionFailed, err)
	}

	p.entries[id] = &poolEntry{conn: conn, lastUsed: time.Now()}
	return conn, nil
}

func (p *Pool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if entry, ok := p.entries[id]; ok {
		_ = entry.conn.Close()
		delete(p.entries, id)
	}
}

func (p *Pool) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		cutoff := time.Now().Add(-ConnTTLMinutes * time.Minute)
		for id, entry := range p.entries {
			if entry.lastUsed.Before(cutoff) {
				_ = entry.conn.Close()
				delete(p.entries, id)
			}
		}
		p.mu.Unlock()
	}
}
