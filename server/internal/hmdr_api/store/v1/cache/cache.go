package cache

import (
	"context"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"github.com/marmotedu/errors"
	"sync"
	"time"
)

type cacheStore struct {
	store          storev1.Factory
	cache          *groupsCache
	expireDuration time.Duration
}

func (c *cacheStore) AgentInfos() storev1.AgentInfoStore {
	return c.store.AgentInfos()
}

func (c *cacheStore) Groups() storev1.GroupStore {
	return newGroupStore(c)
}

func (c *cacheStore) Hosts() storev1.HostStore {
	return newHostStore(c)
}

func (c *cacheStore) WebServerConfigs() storev1.WebServerConfigStore {
	return newWebServerConfigStore(c)
}

func (c *cacheStore) WebServerLogWatchers() storev1.WebServerLogWatcherStore {
	return c.store.WebServerLogWatchers()
}

func (c *cacheStore) WebServerStatistics() storev1.WebServerStatisticsStore {
	return c.store.WebServerStatistics()
}

func (c *cacheStore) WebServerBinCMD() storev1.WebServerBinCMDStore {
	return c.store.WebServerBinCMD()
}

func (c *cacheStore) Close() error {
	return c.store.Close()
}

func (c *cacheStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
	err := c.flushConfigCache(ctx, opts)
	if err != nil {
		return nil, err
	}
	return c.getConfig(opts)
}

func (c *cacheStore) flushConfigCache(ctx context.Context, opts metav1.WebServerOptions) error {
	cache := c.cache.GetGroup(opts.GroupID).GetHost(opts.HostID).GetConfigCache(opts.ServerName)
	if cache.hasExpired(c.expireDuration) {
		cache.release()
		config, err := c.store.WebServerConfigs().GetConfig(ctx, opts)
		if err != nil {
			return err
		}
		return cache.refresh(config)
	}
	return nil
}

func (c *cacheStore) getConfig(opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
	cache := c.cache.GetGroup(opts.GroupID).GetHost(opts.HostID).GetConfigCache(opts.ServerName)
	if !cache.available() {
		return nil, errors.Errorf("nginx config cache of the server(%v) is unavailable", opts)
	}
	return cache.config, nil
}

func (c *cacheStore) CheckBeforeUpdating(ctx context.Context, opts metav1.WebServerOptions) error {
	remoteConfig, err := c.store.WebServerConfigs().GetConfig(ctx, opts)
	if err != nil {
		return errors.Wrapf(err, "failed to check for the server(%v), before updating config", opts)
	}
	if c.cache.GetGroup(opts.GroupID).
		GetHost(opts.HostID).
		GetConfigCache(opts.ServerName).
		isDiff(remoteConfig) {
		return errors.Errorf("failed to check for the server(%v), before updating config. Because the config has been changed", opts)
	}
	return nil
}

type groupsCache struct {
	cache  map[uint]*hostsCache
	locker *sync.RWMutex
}

func (g *groupsCache) GetGroup(groupid uint) *hostsCache {
	g.locker.RLock()
	defer g.locker.RUnlock()
	if _, has := g.cache[groupid]; !has {
		g.cache[groupid] = &hostsCache{
			cache: make(map[uint]*serversCache),
			lock:  new(sync.RWMutex),
		}
	}
	return g.cache[groupid]
}

func (g *groupsCache) ReleaseGroup(groupid uint) *groupsCache {
	g.locker.Lock()
	defer g.locker.Unlock()
	delete(g.cache, groupid)
	return g
}

type hostsCache struct {
	cache map[uint]*serversCache
	lock  *sync.RWMutex
}

func (h *hostsCache) GetHost(hostid uint) *serversCache {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if _, has := h.cache[hostid]; !has {
		h.cache[hostid] = &serversCache{
			cache: make(map[string]*nginxConfigCache),
			lock:  new(sync.RWMutex),
		}
	}
	return h.cache[hostid]
}
func (h *hostsCache) ReleaseHost(hostid uint) *hostsCache {
	h.lock.Lock()
	defer h.lock.RLock()
	delete(h.cache, hostid)
	return h
}

type serversCache struct {
	cache map[string]*nginxConfigCache
	lock  *sync.RWMutex
}

func (s *serversCache) GetConfigCache(servername string) *nginxConfigCache {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if _, has := s.cache[servername]; !has {
		s.cache[servername] = &nginxConfigCache{rwLocker: new(sync.RWMutex)}
	}
	return s.cache[servername]
}

func (s *serversCache) ReleaseConfigCache(servername string) *serversCache {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.cache, servername)
	return s
}

type nginxConfigCache struct {
	config        configuration.NginxConfig
	fingerprinter utilsV3.ConfigFingerprinter

	rwLocker *sync.RWMutex
}

func (c *nginxConfigCache) isDiff(config configuration.NginxConfig) bool {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.fingerprinter == nil || c.fingerprinter.Diff(utilsV3.NewConfigFingerprinter(config.Dump()))
}

func (c *nginxConfigCache) hasExpired(expireDuration time.Duration) bool {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.config == nil || c.config.UpdatedTimestamp().Add(expireDuration).Before(time.Now())
}

func (c *nginxConfigCache) release() {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	c.config = nil
	c.fingerprinter = nil
}

func (c *nginxConfigCache) available() bool {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.config != nil && c.fingerprinter != nil
}

func (c *nginxConfigCache) refresh(config configuration.NginxConfig) error {
	c.release()
	if config == nil {
		return errors.New("refresh into a nil config")
	}
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	c.config = config
	c.fingerprinter = utilsV3.NewConfigFingerprinterWithTimestamp(config.Dump(), config.UpdatedTimestamp())
	return nil
}

func GetCacheStore(store storev1.Factory, cacheExpireDur time.Duration) storev1.Factory {
	return &cacheStore{
		store: store,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: new(sync.RWMutex),
		},
		expireDuration: cacheExpireDur,
	}
}
