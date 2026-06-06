package namer

import (
	"math/rand"
	"strings"
	"testing"
)

const testYAML = `
common:
  star: 1
  rocket: 1
adjective:
  blazing: 1
  silent: 1
substantive:
  wizard: 1
  ninja: 1
intro:
  "Behold:": 1
outro:
  "is unstoppable.": 1
`

func newTestNamer(t *testing.T, seed int64) *Namer {
	t.Helper()
	n, err := New([]byte(testYAML), WithRand(rand.New(rand.NewSource(seed))))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return n
}

func TestNewLoadsYAML(t *testing.T) {
	n := newTestNamer(t, 1)
	if _, ok := n.Words()["common"]; !ok {
		t.Fatal("common list not loaded")
	}
}

func TestNewRejectsInvalidYAML(t *testing.T) {
	if _, err := New([]byte("not: valid: yaml: at: all")); err == nil {
		t.Fatal("expected error from invalid YAML")
	}
}

func TestSuggestReturnsNonEmptyTitleCased(t *testing.T) {
	n := newTestNamer(t, 1)
	name := n.Suggest()
	if name == "" {
		t.Fatal("Suggest returned empty string")
	}
	parts := strings.Fields(name)
	if len(parts) != 2 {
		t.Fatalf("Suggest = %q, want two words", name)
	}
	for _, p := range parts {
		if p[0] < 'A' || p[0] > 'Z' {
			t.Fatalf("Suggest = %q, want title case", name)
		}
	}
}

func TestSuggestNReturnsRequestedCount(t *testing.T) {
	n := newTestNamer(t, 1)
	names := n.SuggestN(5)
	if len(names) != 5 {
		t.Fatalf("len(SuggestN(5)) = %d", len(names))
	}
	for i, name := range names {
		if name == "" {
			t.Fatalf("names[%d] empty", i)
		}
	}
}

func TestSuggestNZero(t *testing.T) {
	n := newTestNamer(t, 1)
	if got := n.SuggestN(0); len(got) != 0 {
		t.Fatalf("SuggestN(0) = %v, want empty", got)
	}
}

func TestSuggestIsDeterministicWithSeed(t *testing.T) {
	a := newTestNamer(t, 42)
	b := newTestNamer(t, 42)
	for i := 0; i < 10; i++ {
		if a.Suggest() != b.Suggest() {
			t.Fatal("same-seed namers diverged")
		}
	}
}
