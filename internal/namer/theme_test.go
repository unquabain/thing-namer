package namer

import (
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func TestRenderContextThemeColorsParseable(t *testing.T) {
	n := newTestNamer(t, 7)
	ctx := n.NewThemedContext()
	colors := []string{
		ctx.Background, ctx.ContentShadow1, ctx.ContentShadow2,
		ctx.Title, ctx.SubTitle, ctx.IntroColor, ctx.OutroColor, ctx.ProjectNameColor,
	}
	for i, c := range colors {
		if _, err := colorful.Hex(c); err != nil {
			t.Fatalf("color %d (%q) not parseable: %v", i, c, err)
		}
	}
}

func TestRenderContextThemeHasAdequateContrast(t *testing.T) {
	n := newTestNamer(t, 7)
	ctx := n.NewThemedContext()
	bg, err := colorful.Hex(ctx.Background)
	if err != nil {
		t.Fatalf("bg: %v", err)
	}
	title, err := colorful.Hex(ctx.Title)
	if err != nil {
		t.Fatalf("title: %v", err)
	}
	// Lab delta-E between bg and title; threshold chosen empirically — the
	// theme is designed so that title sits clearly against the background.
	if d := bg.DistanceLab(title); d < 0.2 {
		t.Fatalf("title/background delta-E = %v, want >= 0.2", d)
	}
}

func TestNewContextHasNameIntroOutro(t *testing.T) {
	n := newTestNamer(t, 7)
	ctx := n.NewContext()
	if ctx.ProjectName == "" {
		t.Fatal("ProjectName empty")
	}
	if ctx.IntroText == "" {
		t.Fatal("IntroText empty")
	}
	if ctx.OutroText == "" {
		t.Fatal("OutroText empty")
	}
}
