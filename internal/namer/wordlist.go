package namer

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/apex/log"
)

type WeightedWord struct {
	Word   string
	Weight int
}

type WordList struct {
	WeightedWords []WeightedWord
}

// Choose picks a word from the list weighted by Weight. The provided
// *rand.Rand is the sole source of randomness.
func (wl *WordList) Choose(r *rand.Rand) string {
	total := 0
	for _, ww := range wl.WeightedWords {
		total += ww.Weight
	}
	if total <= 0 {
		return ""
	}
	pick := r.Intn(total)
	cum := 0
	for _, ww := range wl.WeightedWords {
		cum += ww.Weight
		if pick < cum {
			return ww.Word
		}
	}
	return wl.WeightedWords[len(wl.WeightedWords)-1].Word
}

// Add returns a new WordList containing the entries of wl followed by other,
// without mutating either input.
func (wl *WordList) Add(other *WordList) *WordList {
	sum := &WordList{
		WeightedWords: make([]WeightedWord, 0, len(wl.WeightedWords)+len(other.WeightedWords)),
	}
	sum.WeightedWords = append(sum.WeightedWords, wl.WeightedWords...)
	sum.WeightedWords = append(sum.WeightedWords, other.WeightedWords...)
	return sum
}

func (wl *WordList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]int)
	if err := unmarshal(&m); err != nil {
		return fmt.Errorf("unable to read a wordlist: %w", err)
	}
	for word, val := range m {
		wl.WeightedWords = append(wl.WeightedWords, WeightedWord{Word: word, Weight: val})
	}
	// Sort so the internal order is independent of map iteration, which would
	// otherwise leak non-determinism into seeded RNG output.
	sort.Slice(wl.WeightedWords, func(i, j int) bool {
		return wl.WeightedWords[i].Word < wl.WeightedWords[j].Word
	})
	return nil
}

// WordFile is a map of named WordLists, typically loaded from YAML.
type WordFile map[string]WordList

// Choose concatenates the named lists and picks a single word from the union.
// If any list name is unknown, Choose logs an error and returns "".
func (wf WordFile) Choose(r *rand.Rand, lists ...string) string {
	combined := &WordList{}
	for _, name := range lists {
		nl, ok := wf[name]
		if !ok {
			log.WithField("listname", name).Error("no such word list")
			return ""
		}
		combined = combined.Add(&nl)
	}
	return combined.Choose(r)
}
