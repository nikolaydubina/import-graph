package basiccache

import (
	"fmt"
	"net/url"
	"sync"
)

// GoURLResolver is implementation to be cached
type GoURLResolver interface {
	ResolveGitHubURL(name string) (url.URL, error)
	ResolveGitURL(name string) (url.URL, error)
}

// GoCachedResolver caches GoURLResolver
type GoCachedResolver struct {
	URLResolver GoURLResolver
	Storage     sync.Map
}

// ResolveGitHubURL cached version
func (c *GoCachedResolver) ResolveGitHubURL(name string) (url.URL, error) {
	return c.tryLoad(name, newKeyGitHubURLKey(name), c.URLResolver.ResolveGitHubURL)
}

// ResolveGitURL cached version
func (c *GoCachedResolver) ResolveGitURL(name string) (url.URL, error) {
	return c.tryLoad(name, newKeyGitURLKey(name), c.URLResolver.ResolveGitURL)
}

// tryLoad will load from Cache or invoke f and set to cache and return
func (c *GoCachedResolver) tryLoad(name string, vkey key, f func(name string) (url.URL, error)) (url.URL, error) {
	val, ok := c.Storage.Load(vkey)
	if !ok {
		nVal, err := f(name)
		if err != nil {
			return url.URL{}, fmt.Errorf("can not get GitHubURL: %w", err)
		}
		c.Storage.Store(vkey, nVal)
		return nVal, nil
	}
	ret, ok := val.(url.URL)
	if !ok {
		return url.URL{}, fmt.Errorf("wrong type: %#v", val)
	}
	return ret, nil
}

// key is cache key
type key string

func newKeyGitHubURLKey(name string) key {
	return key("GitHubURL: " + name)
}

func newKeyGitURLKey(name string) key {
	return key("GitURL: " + name)
}
