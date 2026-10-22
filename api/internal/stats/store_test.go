package stats_test

import (
	"sync"
	"testing"

	"github.com/K-T-Ng/ktblog/api/internal/stats"
)

func TestViewNewSlug(t *testing.T) {
	// Arrange
	store := stats.NewStore()

	// Act
	count := store.View("a")

	// Assert
	expected := stats.Counts{Views: 1, Likes: 0}
	if count != expected {
		t.Errorf("View new slug got = %+v, want = %+v", count, expected)
	}
}

func TestViewTwice(t *testing.T) {
	// Arrange
	store := stats.NewStore()

	// Act
	store.View("a")
	count := store.View("a")

	// Assert
	expected := stats.Counts{Views: 2, Likes: 0}
	if count != expected {
		t.Errorf("View new slug twice got = %+v, want = %+v", count, expected)
	}
}

func TestLikeAndUnlike(t *testing.T) {
	// Arrange
	store := stats.NewStore()

	// Act
	store.Like("a")
	store.Like("a")
	count := store.Unlike("a")

	// Assert
	expected := stats.Counts{Views: 0, Likes: 1}
	if count != expected {
		t.Errorf("Like, like, unlike got = %+v, want = %+v", count, expected)
	}
}

func TestUnlikeFloorsAtZero(t *testing.T) {
	// Arrange
	store := stats.NewStore()

	// Act
	count := store.Unlike("a")

	// Assert
	expected := stats.Counts{Views: 0, Likes: 0}
	if count != expected {
		t.Errorf("unlike slug got = %+v, want = %+v", count, expected)
	}
}

func TestSlugsAreIndependent(t *testing.T) {
	// Arrange
	store := stats.NewStore()

	// Act
	store.Like("a")
	countB := store.View("b")
	countA := store.Like("a")

	// Assert
	expectedA := stats.Counts{Views: 0, Likes: 2}
	if countA != expectedA {
		t.Errorf("Slug a got = %+v, want = %+v", countA, expectedA)
	}
	expectedB := stats.Counts{Views: 1, Likes: 0}
	if countB != expectedB {
		t.Errorf("Slug b got = %+v, want = %+v", countB, expectedB)
	}
}

func TestConcurrentViews(t *testing.T) {
	// Arrange
	store := stats.NewStore()

	// Act
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			for range 1000 {
				store.View("a")
			}
		})
	}
	wg.Wait()
	count := store.View("a")

	// Assert
	expected := stats.Counts{Views: 100001, Likes: 0}
	if count != expected {
		t.Errorf("Concurrent views got = %+v, want = %+v", count, expected)
	}
}
