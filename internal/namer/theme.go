package namer

import (
	"github.com/lucasb-eyer/go-colorful"
)

// RenderContext is the data passed to HTTP templates (HTML/JSON/Go).
type RenderContext struct {
	ProjectName       string `json:"projectName"`
	Background        string `json:"-"`
	ContentBackground string `json:"-"`
	ContentShadow1    string `json:"-"`
	ContentShadow2    string `json:"-"`
	Title             string `json:"-"`
	SubTitle          string `json:"-"`
	IntroText         string `json:"intro"`
	IntroColor        string `json:"-"`
	OutroText         string `json:"outro"`
	OutroColor        string `json:"-"`
	ProjectNameColor  string `json:"-"`
}

var (
	white, _ = colorful.Hex("#FFFFFF")
	black, _ = colorful.Hex("#000000")
)

// NewContext returns a RenderContext with just the textual fields filled in.
func (n *Namer) NewContext() RenderContext {
	return RenderContext{
		ProjectName: n.Suggest(),
		IntroText:   n.ChooseIntro(),
		OutroText:   n.ChooseOutro(),
	}
}

// NewThemedContext returns a RenderContext with both text and a generated
// color theme.
func (n *Namer) NewThemedContext() RenderContext {
	ctx := n.NewContext()
	n.applyTheme(&ctx)
	return ctx
}

func (n *Namer) applyTheme(ctx *RenderContext) {
	h, c, l := n.rng.Float64()*360, 0.3, 0.7
	main := colorful.Hcl(h, c, l)
	complement := colorful.Hcl(h+180.0, 0.6, 0.4)
	ctx.Background = main.BlendLab(white, 0.75).Clamped().Hex()
	ctx.ContentBackground = "white"
	ctx.ContentShadow1 = main.BlendLab(white, 0.5).Clamped().Hex()
	ctx.ContentShadow2 = main.Clamped().Hex()
	ctx.Title = main.BlendLab(black, 0.2).Clamped().Hex()
	ctx.SubTitle = ctx.ContentShadow2
	ctx.IntroColor = complement.BlendLab(black, 0.6).Clamped().Hex()
	ctx.OutroColor = ctx.Title
	ctx.ProjectNameColor = complement.Clamped().Hex()
}
