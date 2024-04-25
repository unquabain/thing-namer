package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

//go:embed data/words.yaml
var words []byte

//go:embed templates/index.html
var indexRaw string
var index = template.Must(template.New(`index`).Funcs(sprig.TxtFuncMap()).Parse(indexRaw))

//go:embed templates/client.go.tmpl
var goRaw string
var goTmplt = template.Must(template.New(`go`).Funcs(sprig.TxtFuncMap()).Parse(goRaw))

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

var white, _ = colorful.Hex("#FFFFFF")
var black, _ = colorful.Hex("#000000")

var defaultContext = RenderContext{
	Background:        `#EEEEEE`,
	ContentBackground: `white`,
	ContentShadow1:    `silver`,
	ContentShadow2:    `silver`,
	Title:             `black`,
	SubTitle:          `black`,
	IntroText:         `Your project is now called`,
	IntroColor:        `black`,
	OutroText:         `You're welcome!`,
	OutroColor:        `black`,
	ProjectNameColor:  `black`,
}

func createTheme(context *RenderContext) {
	h, c, l := rand.Float64()*360, 0.3, 0.7
	main := colorful.Hcl(h, c, l)
	complement := colorful.Hcl(h+180.0, 0.6, 0.4)
	context.Background = main.BlendLab(white, 0.75).Clamped().Hex()
	context.ContentBackground = `white`
	context.ContentShadow1 = main.BlendLab(white, 0.5).Clamped().Hex()
	context.ContentShadow2 = main.Clamped().Hex()
	context.Title = main.BlendLab(black, 0.2).Clamped().Hex()
	context.SubTitle = context.ContentShadow2
	context.IntroColor = complement.BlendLab(black, 0.6).Clamped().Hex()
	context.OutroColor = context.Title
	context.ProjectNameColor = complement.Clamped().Hex()
}

func (wf WordFile) createContext() RenderContext {
	return RenderContext{
		ProjectName: wf.projectName(),
		IntroText:   wf.Choose(`intro`),
		OutroText:   wf.Choose(`outro`),
	}
}

func (wf WordFile) createThemedContext() RenderContext {
	context := wf.createContext()
	createTheme(&context)
	return context
}

func (wf WordFile) projectName() string {
	var adjective,
		substantive string
	adjective = wf.Choose(`common`, `adjective`)
	for {
		substantive = wf.Choose(`common`, `substantive`)
		if substantive != adjective {
			break
		}
	}
	title := cases.Title(language.English)
	return title.String(fmt.Sprintf(`%s %s`, adjective, substantive))
}

func requestIsJSON(r *http.Request) bool {
	if strings.HasSuffix(r.URL.Path, `.json`) {
		return true
	}
	if r.Header.Get(`Accept`) == `application/json` {
		return true
	}
	return false
}

func requestIsGo(r *http.Request) bool {
	if strings.HasSuffix(r.URL.Path, `.go`) {
		return true
	}
	return false
}

func (WordFile) renderGo(w http.ResponseWriter, r *http.Request) {
	proto := `http`
	if r.TLS != nil {
		proto = `https`
	}
	var renderContext = struct {
		Server string
	}{fmt.Sprintf(`%s://%s`, proto, r.Host)}
	w.Header().Add(`Content-Type`, `text/plain`)
	w.Header().Add(`Content-Disposition`, `attachment; filename="wizardbacon.go"`)
	goTmplt.Execute(w, renderContext)
}

func (wf WordFile) renderHTML(w http.ResponseWriter) {
	w.Header().Add(`Content-Type`, `text/html`)
	context := wf.createThemedContext()
	if err := index.Execute(w, context); err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err.Error())
		return
	}
}

func (wf WordFile) renderJSON(w http.ResponseWriter) {

	w.Header().Add(`Content-Type`, `application/json`)
	context := wf.createContext()
	if err := json.NewEncoder(w).Encode(context); err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err.Error())
		return
	}
}

func (wf WordFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if requestIsJSON(r) {
		wf.renderJSON(w)
	} else if requestIsGo(r) {
		wf.renderGo(w, r)
	} else {
		wf.renderHTML(w)
	}
}

func main() {
	wf := make(WordFile)

	yaml.Unmarshal(words, wf)
	http.ListenAndServe(`:9099`, wf)

}
