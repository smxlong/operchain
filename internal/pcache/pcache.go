package pcache

import "sync"

// Predicate represents a cacheable boolean function.
type Predicate struct {
	f func() bool
}

// Cache is a predicate value Cache.
type Cache struct {
	c    map[*Predicate]bool
	lock sync.Mutex
}

// New creates a new Cache.
func New() *Cache {
	return &Cache{
		c: map[*Predicate]bool{},
	}
}

// Eval evaluates the predicate in the cache.
func (c *Cache) Eval(p *Predicate) bool {
	if val, ok := c.isInCache(p); ok {
		return val
	}
	val := p.f()
	c.addToCache(p, val)
	return val
}

// And returns a new Predicate that is the logical AND of the given Predicates.
func (c *Cache) And(p ...*Predicate) *Predicate {
	return &Predicate{
		f: func() bool {
			for _, pred := range p {
				if !c.Eval(pred) {
					return false
				}
			}
			return true
		},
	}
}

// Or returns a new Predicate that is the logical OR of the given Predicates.
func (c *Cache) Or(p ...*Predicate) *Predicate {
	return &Predicate{
		f: func() bool {
			for _, pred := range p {
				if c.Eval(pred) {
					return true
				}
			}
			return false
		},
	}
}

// isInCache returns true if the given predicate is in the cache.
func (c *Cache) isInCache(p *Predicate) (bool, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := c.c[p]
	return val, ok
}

// addToCache adds the given value to the cache for the given predicate.
func (c *Cache) addToCache(p *Predicate, value bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.c[p] = value
}
