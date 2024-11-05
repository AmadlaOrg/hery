package cache

// NewCacheService to set up the entity Cache service
func NewCacheService() ICache {
	return &SCache{}
}
