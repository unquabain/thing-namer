package namer

import (
	"math/rand"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestWordListChooseReturnsAWordFromTheList(t *testing.T) {
	wl := &WordList{WeightedWords: []WeightedWord{
		{Word: "alpha", Weight: 1},
		{Word: "beta", Weight: 1},
		{Word: "gamma", Weight: 1},
	}}
	r := rand.New(rand.NewSource(42))
	got := wl.Choose(r)
	if got != "alpha" && got != "beta" && got != "gamma" {
		t.Fatalf("Choose returned %q, not in list", got)
	}
}

func TestWordListChooseRespectsWeights(t *testing.T) {
	wl := &WordList{WeightedWords: []WeightedWord{
		{Word: "rare", Weight: 1},
		{Word: "common", Weight: 99},
	}}
	r := rand.New(rand.NewSource(1))
	counts := map[string]int{}
	const N = 10000
	for i := 0; i < N; i++ {
		counts[wl.Choose(r)]++
	}
	if counts["common"] <= 2*counts["rare"] {
		t.Fatalf("expected common >> rare, got common=%d rare=%d", counts["common"], counts["rare"])
	}
}

func TestWordListAddConcatenatesWithoutMutating(t *testing.T) {
	a := &WordList{WeightedWords: []WeightedWord{{Word: "a", Weight: 1}}}
	b := &WordList{WeightedWords: []WeightedWord{{Word: "b", Weight: 2}}}
	sum := a.Add(b)
	if len(a.WeightedWords) != 1 || a.WeightedWords[0].Word != "a" {
		t.Fatalf("a was mutated: %+v", a.WeightedWords)
	}
	if len(b.WeightedWords) != 1 || b.WeightedWords[0].Word != "b" {
		t.Fatalf("b was mutated: %+v", b.WeightedWords)
	}
	if len(sum.WeightedWords) != 2 {
		t.Fatalf("sum length = %d, want 2", len(sum.WeightedWords))
	}
}

func TestWordListUnmarshalYAML(t *testing.T) {
	doc := []byte("foo: 3\nbar: 7\n")
	var wl WordList
	if err := yaml.Unmarshal(doc, &wl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(wl.WeightedWords) != 2 {
		t.Fatalf("len = %d, want 2", len(wl.WeightedWords))
	}
	seen := map[string]int{}
	for _, ww := range wl.WeightedWords {
		seen[ww.Word] = ww.Weight
	}
	if seen["foo"] != 3 || seen["bar"] != 7 {
		t.Fatalf("got %+v, want foo=3 bar=7", seen)
	}
}

func TestWordFileChooseAcrossMultipleLists(t *testing.T) {
	wf := WordFile{
		"a": WordList{WeightedWords: []WeightedWord{{Word: "x", Weight: 1}}},
		"b": WordList{WeightedWords: []WeightedWord{{Word: "y", Weight: 1}}},
	}
	r := rand.New(rand.NewSource(42))
	got := wf.Choose(r, "a", "b")
	if got != "x" && got != "y" {
		t.Fatalf("Choose returned %q, expected x or y", got)
	}
}

func TestWordFileChooseUnknownListReturnsEmpty(t *testing.T) {
	wf := WordFile{}
	r := rand.New(rand.NewSource(1))
	if got := wf.Choose(r, "nope"); got != "" {
		t.Fatalf("Choose unknown list = %q, want empty", got)
	}
}
