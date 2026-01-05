package clients

import (
	"context"
	"sync"
	"time"

	"ms-parcel-core/internal/parcel/parcel_core/port"
)

type tenantOptionsCacheEntry struct {
	opts      port.ParcelOptions
	expiresAt time.Time
}

type CachedTenantOptionsProvider struct {
	inner port.TenantOptionsProvider
	ttl   time.Duration

	mu    sync.Mutex
	cache map[string]tenantOptionsCacheEntry
}

var _ port.TenantOptionsProvider = (*CachedTenantOptionsProvider)(nil)

func NewCachedTenantOptionsProvider(inner port.TenantOptionsProvider, ttl time.Duration) *CachedTenantOptionsProvider {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	return &CachedTenantOptionsProvider{inner: inner, ttl: ttl, cache: map[string]tenantOptionsCacheEntry{}}
}

func (p *CachedTenantOptionsProvider) GetParcelOptions(ctx context.Context, tenantID string) (port.ParcelOptions, error) {
	now := time.Now().UTC()

	p.mu.Lock()
	if p.cache != nil {
		if e, ok := p.cache[tenantID]; ok {
			if now.Before(e.expiresAt) {
				opts := e.opts
				p.mu.Unlock()
				return opts, nil
			}
		}
	}
	p.mu.Unlock()

	opts, err := p.inner.GetParcelOptions(ctx, tenantID)
	if err != nil {
		return port.ParcelOptions{}, err
	}

	p.mu.Lock()
	if p.cache == nil {
		p.cache = map[string]tenantOptionsCacheEntry{}
	}
	p.cache[tenantID] = tenantOptionsCacheEntry{opts: opts, expiresAt: now.Add(p.ttl)}
	p.mu.Unlock()

	return opts, nil
}
