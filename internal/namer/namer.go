package namer

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

// Namer generates project names from a loaded WordFile.
type Namer struct {
	words WordFile
	rng   *rand.Rand
}

// Option configures a Namer at construction time.
type Option func(*Namer)

// WithRand replaces the default time-seeded RNG with the provided one.
// Tests use this to get deterministic output.
func WithRand(r *rand.Rand) Option {
	return func(n *Namer) { n.rng = r }
}

// New parses the YAML word file and returns a ready-to-use *Namer.
// By default the RNG is seeded with the current time; use WithRand to override.
func New(yamlBytes []byte, opts ...Option) (*Namer, error) {
	wf := make(WordFile)
	if err := yaml.Unmarshal(yamlBytes, wf); err != nil {
		return nil, fmt.Errorf("namer: parse YAML: %w", err)
	}
	n := &Namer{
		words: wf,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, opt := range opts {
		opt(n)
	}
	return n, nil
}

// Words exposes the loaded WordFile for callers (e.g. theme rendering)
// that need direct access to named lists.
func (n *Namer) Words() WordFile { return n.words }

// Rand returns the Namer's RNG. Callers that need additional randomness
// (e.g. theme generation) should use this rather than math/rand globals.
func (n *Namer) Rand() *rand.Rand { return n.rng }

// Suggest returns one project name: a title-cased "adjective substantive"
// pair drawn from the common+adjective and common+substantive lists.
// It re-rolls until the adjective and substantive differ.
func (n *Namer) Suggest() string {
	adjective := n.words.Choose(n.rng, "common", "adjective")
	var substantive string
	for {
		substantive = n.words.Choose(n.rng, "common", "substantive")
		if substantive != adjective {
			break
		}
	}
	title := cases.Title(language.English)
	return title.String(fmt.Sprintf("%s %s", adjective, substantive))
}

// SuggestN returns count independent suggestions. Duplicates are allowed.
func (n *Namer) SuggestN(count int) []string {
	if count <= 0 {
		return []string{}
	}
	out := make([]string, count)
	for i := 0; i < count; i++ {
		out[i] = n.Suggest()
	}
	return out
}

// ChooseIntro and ChooseOutro pull from the intro/outro lists; used by the
// HTTP render layer.
func (n *Namer) ChooseIntro() string { return n.words.Choose(n.rng, "intro") }
func (n *Namer) ChooseOutro() string { return n.words.Choose(n.rng, "outro") }
