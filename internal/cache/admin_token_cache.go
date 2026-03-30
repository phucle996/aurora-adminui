package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"aurora-adminui/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

const (
	runtimeTokenPrefix  = "admin:token:"
	bootstrapCurrentKey = "admin:bootstrap:current"
	preauthPrefix       = "admin:preauth:"
	sessionPrefix       = "admin:session:"
	totpSetupPrefix     = "admin:totp-setup:"
)

type cacheEntry struct {
	tokenType entity.TokenType
	expiresAt time.Time
}

type AdminTokenCache struct {
	redis *redis.Client
	mu    sync.RWMutex
	items map[string]cacheEntry
}

func NewAdminTokenCache(redisClient *redis.Client) *AdminTokenCache {
	return &AdminTokenCache{redis: redisClient, items: make(map[string]cacheEntry)}
}

func (c *AdminTokenCache) GetTokenType(ctx context.Context, token string) (entity.TokenType, bool, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false, nil
	}
	if tokenType, ok := c.getMemory(token); ok {
		return tokenType, true, nil
	}
	if c == nil || c.redis == nil {
		return "", false, nil
	}
	val, err := c.redis.Get(ctx, runtimeTokenPrefix+token).Result()
	if err == nil {
		tokenType := entity.TokenType(strings.TrimSpace(val))
		if tokenType == entity.TokenTypeBootstrap {
			c.setMemory(token, tokenType, 15*time.Minute)
		} else {
			c.setMemory(token, tokenType, 5*time.Minute)
		}
		return tokenType, true, nil
	}
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	return "", false, err
}

func (c *AdminTokenCache) SetToken(ctx context.Context, token string, tokenType entity.TokenType, ttl time.Duration) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("token is empty")
	}
	memoryTTL := ttl
	if tokenType == entity.TokenTypeAPIToken {
		memoryTTL = 5 * time.Minute
	}
	c.setMemory(token, tokenType, memoryTTL)
	if c == nil || c.redis == nil {
		return nil
	}
	return c.redis.Set(ctx, runtimeTokenPrefix+token, string(tokenType), ttl).Err()
}

func (c *AdminTokenCache) DeleteToken(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}
	c.mu.Lock()
	delete(c.items, token)
	c.mu.Unlock()
	if c == nil || c.redis == nil {
		return nil
	}
	return c.redis.Del(ctx, runtimeTokenPrefix+token).Err()
}

func (c *AdminTokenCache) SetBootstrapPointer(ctx context.Context, token string, ttl time.Duration) error {
	if c == nil || c.redis == nil {
		return nil
	}
	return c.redis.Set(ctx, bootstrapCurrentKey, strings.TrimSpace(token), ttl).Err()
}

func (c *AdminTokenCache) GetBootstrapPointer(ctx context.Context) (string, error) {
	if c == nil || c.redis == nil {
		return "", nil
	}
	val, err := c.redis.Get(ctx, bootstrapCurrentKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}

func (c *AdminTokenCache) ClearBootstrapPointer(ctx context.Context) error {
	if c == nil || c.redis == nil {
		return nil
	}
	return c.redis.Del(ctx, bootstrapCurrentKey).Err()
}

func (c *AdminTokenCache) SetPreauth(ctx context.Context, preauth *entity.AdminPreauthSession, ttl time.Duration) error {
	if c == nil || c.redis == nil {
		return errors.New("redis client is nil")
	}
	if preauth == nil || strings.TrimSpace(preauth.ID) == "" {
		return errors.New("preauth session is nil")
	}
	payload, err := json.Marshal(preauth)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, preauthPrefix+preauth.ID, payload, ttl).Err()
}

func (c *AdminTokenCache) GetPreauth(ctx context.Context, id string) (*entity.AdminPreauthSession, error) {
	if c == nil || c.redis == nil {
		return nil, errors.New("redis client is nil")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("preauth session id is empty")
	}
	raw, err := c.redis.Get(ctx, preauthPrefix+id).Bytes()
	if err != nil {
		return nil, err
	}
	var preauth entity.AdminPreauthSession
	if err := json.Unmarshal(raw, &preauth); err != nil {
		return nil, err
	}
	return &preauth, nil
}

func (c *AdminTokenCache) DeletePreauth(ctx context.Context, id string) error {
	if c == nil || c.redis == nil {
		return errors.New("redis client is nil")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	return c.redis.Del(ctx, preauthPrefix+id).Err()
}

func (c *AdminTokenCache) SetSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	if c == nil || c.redis == nil {
		return errors.New("redis client is nil")
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return errors.New("session id is empty")
	}
	return c.redis.Set(ctx, sessionPrefix+sessionID, "active", ttl).Err()
}

func (c *AdminTokenCache) HasSession(ctx context.Context, sessionID string) (bool, error) {
	if c == nil || c.redis == nil {
		return false, errors.New("redis client is nil")
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return false, nil
	}
	count, err := c.redis.Exists(ctx, sessionPrefix+sessionID).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *AdminTokenCache) DeleteSession(ctx context.Context, sessionID string) error {
	if c == nil || c.redis == nil {
		return errors.New("redis client is nil")
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil
	}
	return c.redis.Del(ctx, sessionPrefix+sessionID).Err()
}

func (c *AdminTokenCache) SetTOTPSetup(ctx context.Context, session *entity.AdminTOTPSetupSession, ttl time.Duration) error {
	if c == nil || c.redis == nil {
		return errors.New("redis client is nil")
	}
	if session == nil || strings.TrimSpace(session.ID) == "" {
		return errors.New("totp setup session is nil")
	}
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, totpSetupPrefix+session.ID, payload, ttl).Err()
}

func (c *AdminTokenCache) GetTOTPSetup(ctx context.Context, id string) (*entity.AdminTOTPSetupSession, error) {
	if c == nil || c.redis == nil {
		return nil, errors.New("redis client is nil")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("totp setup session id is empty")
	}
	raw, err := c.redis.Get(ctx, totpSetupPrefix+id).Bytes()
	if err != nil {
		return nil, err
	}
	var setup entity.AdminTOTPSetupSession
	if err := json.Unmarshal(raw, &setup); err != nil {
		return nil, err
	}
	return &setup, nil
}

func (c *AdminTokenCache) DeleteTOTPSetup(ctx context.Context, id string) error {
	if c == nil || c.redis == nil {
		return errors.New("redis client is nil")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	return c.redis.Del(ctx, totpSetupPrefix+id).Err()
}

func (c *AdminTokenCache) getMemory(token string) (entity.TokenType, bool) {
	if c == nil {
		return "", false
	}
	c.mu.RLock()
	entry, ok := c.items[token]
	c.mu.RUnlock()
	if !ok {
		return "", false
	}
	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.items, token)
		c.mu.Unlock()
		return "", false
	}
	return entry.tokenType, true
}

func (c *AdminTokenCache) setMemory(token string, tokenType entity.TokenType, ttl time.Duration) {
	if c == nil {
		return
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	c.mu.Lock()
	c.items[token] = cacheEntry{tokenType: tokenType, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}
