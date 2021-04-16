package resolver_cached

import (
	"fmt"
	"net/url"
	"sync"
)

type GoResolver interface {
	ResolveGitHubURL(name string) (*url.URL, error)
	ResolveGitURL(name string) (*url.URL, error)
}

type GoCachedResolver struct {
	Resolver GoResolver
	Storage  sync.Map
}

func (c *GoCachedResolver) ResolveGitHubURL(name string) (*url.URL, error) {
	return c.tryLoad(name, newKeyGitHubURLKey(name), c.Resolver.ResolveGitHubURL)
}

func (c *GoCachedResolver) ResolveGitURL(name string) (*url.URL, error) {
	return c.tryLoad(name, newKeyGitURLKey(name), c.Resolver.ResolveGitURL)
}

// tryLoad will load from Cache or invoke f and set to cache and return
func (c *GoCachedResolver) tryLoad(name string, vkey key, f func(name string) (*url.URL, error)) (*url.URL, error) {
	val, ok := c.Storage.Load(vkey)
	if !ok {
		nVal, err := f(name)
		if err != nil {
			return nil, fmt.Errorf("can not get GitHubURL: %w", err)
		}
		c.Storage.Store(vkey, nVal)
		return nVal, nil
	}
	ret, ok := val.(*url.URL)
	if !ok {
		return nil, fmt.Errorf("wrong type: %#v", val)
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
