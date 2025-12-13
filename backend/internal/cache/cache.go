package cache

import (
	"sync"
	"time"
)

type ComparisonCache struct {
	mu          sync.RWMutex
	comparisons map[int64]*CachedComparison
}

type CachedComparison struct {
	ComparisonID int64
	SetAID       int64
	SetBID       int64
	VotesA       int
	VotesB       int
	ExpiresAt    time.Time
}

// LeaderboardCache can be added later if needed
// For now, we'll use direct queries or add caching in the handler
// This keeps the stack simple - can add Redis later if query performance becomes an issue

func NewComparisonCache() *ComparisonCache {
	return &ComparisonCache{
		comparisons: make(map[int64]*CachedComparison),
	}
}

func (c *ComparisonCache) Get(comparisonID int64) (*CachedComparison, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	comp, exists := c.comparisons[comparisonID]
	if !exists {
		return nil, false
	}

	if time.Now().After(comp.ExpiresAt) {
		delete(c.comparisons, comparisonID)
		return nil, false
	}

	return comp, true
}

func (c *ComparisonCache) Set(comparisonID int64, setAID, setBID int64, votesA, votesB int, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.comparisons[comparisonID] = &CachedComparison{
		ComparisonID: comparisonID,
		SetAID:       setAID,
		SetBID:       setBID,
		VotesA:       votesA,
		VotesB:       votesB,
		ExpiresAt:    time.Now().Add(ttl),
	}
}

func (c *ComparisonCache) UpdateVotes(comparisonID int64, votesA, votesB int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if comp, exists := c.comparisons[comparisonID]; exists {
		comp.VotesA = votesA
		comp.VotesB = votesB
	}
}

func (c *ComparisonCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for id, comp := range c.comparisons {
		if now.After(comp.ExpiresAt) {
			delete(c.comparisons, id)
		}
	}
}

