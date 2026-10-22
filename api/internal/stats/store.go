// Package stats counts post views and likes in memory.
package stats

import "sync"

// Counts holds the view and like counts of a post.
type Counts struct {
	Views int64
	Likes int64
}

// Store keeps per-post view and like counts in memory. It is safe for concurrent use.
type Store struct {
	mu     sync.Mutex
	counts map[string]Counts
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{counts: make(map[string]Counts)}
}

// View increments the view count of slug and returns the updated counts.
func (s *Store) View(slug string) Counts {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := s.counts[slug]
	c.Views++
	s.counts[slug] = c
	return c
}

// Like increments the like count of slug and returns the updated counts.
func (s *Store) Like(slug string) Counts {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := s.counts[slug]
	c.Likes++
	s.counts[slug] = c
	return c
}

// Unlike decrements the like count of slug, never going below zero, and returns the updated counts.
func (s *Store) Unlike(slug string) Counts {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := s.counts[slug]
	if c.Likes > 0 {
		c.Likes--
	}
	s.counts[slug] = c
	return c
}
