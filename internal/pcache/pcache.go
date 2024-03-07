package pcache

import "sync"

// Predicate represents a cacheable boolean function.
type Predicate struct {
	f func(c *Cache) bool
}

// NewPredicate creates a new Predicate.
func NewPredicate(f func() bool) *Predicate {
	return &Predicate{
		f: func(c *Cache) bool {
			return f()
		},
	}
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
func (p *Predicate) Eval(c *Cache) bool {
	if val, ok := c.isInCache(p); ok {
		return val
	}
	val := p.f(c)
	c.addToCache(p, val)
	return val
}

// And returns a new Predicate that is the logical AND of the given Predicates.
func And(p ...*Predicate) *Predicate {
	return &Predicate{
		f: func(c *Cache) bool {
			for _, expr := range p {
				if !expr.Eval(c) {
					return false
				}
			}
			return true
		},
	}
}

// Or returns a new Predicate that is the logical OR of the given Predicates.
func Or(p ...*Predicate) *Predicate {
	return &Predicate{
		f: func(c *Cache) bool {
			for _, expr := range p {
				if expr.Eval(c) {
					return true
				}
			}
			return false
		},
	}
}

// Not returns the negation of the given Predicate.
func Not(p *Predicate) *Predicate {
	return &Predicate{
		f: func(c *Cache) bool {
			return !p.Eval(c)
		},
	}
}

// True returns a Predicate that always returns true.
func True() *Predicate {
	return &Predicate{
		f: func(*Cache) bool {
			return true
		},
	}
}

// False returns a Predicate that always returns false.
func False() *Predicate {
	return &Predicate{
		f: func(*Cache) bool {
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
