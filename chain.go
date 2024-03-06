package operchain

import (
	"context"
	"sync"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/smxlong/operchain/internal/pcache"
)

// Chain is a chain of operchain Rules.
type Chain struct {
	// Rules is the list of rules in the chain.
	Rules []Rule
	// Resources are the resources to load before running the chain.
	Resources interface{}

	lock     sync.Mutex
	stop     bool
	err      error
	interval time.Duration
}

// Action is an action to take in an operchain.
type Action func(context.Context)

// predicate is a private alias for the pcache predicate type to hide it from
// the public API.
type predicate = pcache.Predicate

// Rule is a rule for the operchain.
type Rule struct {
	// When is the predicate for the rule.
	When *predicate
	// Do is the action to take when the predicate is true.
	Do Action
}

// Predicate returns a predicate for the given function.
func Predicate(f func() bool) *predicate {
	return pcache.NewPredicate(f)
}

// Run runs an operchain.
func (c *Chain) Run(ctx context.Context, cache *pcache.Cache) (ctrl.Result, error) {
	c.stop = false
	c.err = nil
	c.interval = 0
	for _, rule := range c.Rules {
		if rule.When == nil || cache.Eval(rule.When) {
			rule.Do(ctx)
			if c.stop || c.err != nil {
				return ctrl.Result{Requeue: true, RequeueAfter: c.interval}, c.err
			}
		}
	}
	return ctrl.Result{Requeue: true, RequeueAfter: c.interval}, nil
}

// Requeue returns an action to set the requeue interval, if it is less than the
// current requeue interval.
func (c *Chain) Requeue(interval time.Duration) Action {
	return func(ctx context.Context) {
		c.doRequeue(interval)
	}
}

func (c *Chain) doRequeue(interval time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if interval > 0 && (c.interval == 0 || interval < c.interval) {
		c.interval = interval
	}
}

// Stop returns an action to stop the operchain.
func (c *Chain) Stop() Action {
	return func(ctx context.Context) {
		c.doStop()
	}
}

func (c *Chain) doStop() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.stop = true
}

// Error returns an action to set the error for the operchain.
func (c *Chain) Error(err error) Action {
	return func(ctx context.Context) {
		c.doError(err)
	}
}

func (c *Chain) doError(err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.err = err
}

// Sequential returns an action that runs the given actions in sequence.
func Sequential(fns ...Action) Action {
	return func(ctx context.Context) {
		for _, fn := range fns {
			fn(ctx)
		}
	}
}

// Parallel returns an action that runs the given actions in parallel.
func Parallel(fns ...Action) Action {
	return func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(len(fns))
		for _, fn := range fns {
			go func(fn Action) {
				defer wg.Done()
				fn(ctx)
			}(fn)
		}
		wg.Wait()
	}
}

// Subchain returns an action that runs the given chain. Any requeue or error
// actions in the subchain will be propagated to the parent chain.
func (c *Chain) Subchain(sub *Chain) Action {
	return func(ctx context.Context) {
		result, err := sub.Run(ctx, nil)
		if err != nil {
			c.doError(err)
		}
		if result.RequeueAfter > 0 {
			c.doRequeue(result.RequeueAfter)
		}
	}
}

// Initialize initializes a chain with the given resources and rules.
func (c *Chain) InitializeChain(resources interface{}, rules []Rule) {
	c.Resources = resources
	c.Rules = rules
}
